package main

import (
        "bufio"
        "fmt"
        "github.com/fgrehm/go-dockerpty"
        "github.com/fsouza/go-dockerclient"
        "os"
)

func (c *container_struct) runContainer() {
        fmt.Println("starting the container:", c.emp_name)
        employes_container_id := c.container.ID
        c.client.StartContainer(employes_container_id, &docker.HostConfig{})
}

func (c *container_struct) createStartConatiner() {
        fmt.Println("creating a container for you:", c.emp_name)
        container, err := c.client.CreateContainer(docker.CreateContainerOptions{
                Name: c.emp_name,
                Config: &docker.Config{
                        Image: c.basic_image,
                        Cmd:   []string{"tail", "-f", "/dev/null"},
                },
        })
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        c.container.ID = container.ID
        c.runContainer()
}

func (c *container_struct) GetStatus() string {
        fmt.Println("getting status:", c.emp_name)
        if c.container.Status != "" {
                return c.container.Status
        } else {
                return "notfound"
        }
}

func (c *container_struct) getTerminal() {
        fmt.Println("getting terminal (to exit type exit):", c.emp_name)
        employes_container_id := c.container.ID
        if employes_container_id == "" {
                panic("employee's container not found")
                os.Exit(1)
        }
        exec, err := c.client.CreateExec(docker.CreateExecOptions{
                Container:    employes_container_id,
                AttachStdin:  true,
                AttachStdout: true,
                AttachStderr: true,
                Tty:          true,
                Cmd:          []string{"/bin/sh"},
        })
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        if err = dockerpty.StartExec(c.client, exec); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        fmt.Println("hope you enjoyed your stay")
}

func (c *container_struct) saveImage() {
        fmt.Println("saving your image with name:", c.emp_name)
        employes_container_id := c.container.ID
        err := c.client.RemoveImage(c.emp_name)
        if err != nil {
                fmt.Println("please Ignore this:", err)
        }
        img, err := c.client.CommitContainer(docker.CommitContainerOptions{
                Container:  employes_container_id,
                Author:     "docker-manager",
                Repository: c.emp_name,
                Tag:        "latest",
        })
        if err != nil {
                fmt.Println("Oh no we couldn't save you image", err)
                os.Exit(1)
        }
        c.image_id = img.ID
        fmt.Println("image saved:", img.ID)
}

type container_struct struct {
        client      *docker.Client
        emp_name    string
        image_id    string
        container   docker.APIContainers
        basic_image string
}

func NewDocker(emp_name string) (*container_struct, error) {
        endpoint := "unix:///var/run/docker.sock"
        var person_container docker.APIContainers
        client, err := docker.NewClient(endpoint)
        if err != nil {
                panic(err)
                return nil, err
        }
        containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
        if err != nil {
                panic(err)
                return nil, err
        }
        for _, container := range containers {
                if name := container.Names[0][1:]; emp_name == name {
                        person_container = container
                }
        }
        container_of_employ := &container_struct{
                client:      client,
                emp_name:    emp_name,
                image_id:    emp_name,
                container:   person_container,
                basic_image: "ubuntu",
        }
        _, err = client.InspectImage(emp_name)
        if err == nil {
                container_of_employ.basic_image = emp_name
        } else {
                fmt.Println("no base image found! if your conatiner is not running we will start with image:", container_of_employ.basic_image)
        }
        return container_of_employ, nil
}

func (c *container_struct) printobject() {
        fmt.Println("client:", c.client)
        fmt.Println("emp_name:", c.emp_name)
        fmt.Println("image_id:", c.image_id)
        fmt.Println("containerID:", c.container.ID)
        fmt.Println("containerstatus:", c.container.Status)
}

func main() {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter employid: ")
        raw_input, _ := reader.ReadString('\n')
        emp_name := raw_input[:len(raw_input)-1]
        docker, err := NewDocker(emp_name)
        if err != nil {
                panic(err)
        }
        docker.printobject()
        fmt.Println("trying to fetch your docker container")
        status := docker.GetStatus()
        fmt.Println("status container:", status)
        if status == "running" {
                docker.getTerminal()
        } else if status == "stopped" {
                docker.runContainer()
                docker.getTerminal()
        } else if status == "exited" || status == "Created" {
                docker.runContainer()
                docker.getTerminal()
        } else if status == "notfound" {
                docker.createStartConatiner()
                docker.getTerminal()
        }
        docker.saveImage()
}
