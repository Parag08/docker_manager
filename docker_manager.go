package main

import (
        "bufio"
        "fmt"
        "github.com/fgrehm/go-dockerpty"
        "github.com/fsouza/go-dockerclient"
        "os"
)

func RunContainer(client *docker.Client, emp_name string) {
        fmt.Println("starting the container:", emp_name)
        employes_container_id := FindContainerID(client, emp_name)
        client.StartContainer(employes_container_id, &docker.HostConfig{})
}

func CreateStartConatiner(client *docker.Client, emp_name string) {
        fmt.Println("creating a container for you:", emp_name)
        _, err := client.CreateContainer(docker.CreateContainerOptions{
                Name: emp_name,
                Config: &docker.Config{
                        Image: "ubuntu",
                        Cmd:   []string{"tail", "-f", "/dev/null"},
                },
        })
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        RunContainer(client, emp_name)
}

func GetStatus(client *docker.Client, emp_name string) string {
        fmt.Println("getting status:", emp_name)
        containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
        if err != nil {
                panic(err)
        }
        for _, container := range containers {
                if name := container.Names[0][1:]; emp_name == name {
                        return container.State
                }
        }
        return "notfound"

}

func FindContainerID(client *docker.Client, emp_name string) string {
        containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
        if err != nil {
                panic(err)
        }
        for _, container := range containers {
                if name := container.Names[0][1:]; emp_name == name {
                        return container.ID
                }
        }
        return ""

}

func GetTerminal(client *docker.Client, emp_name string) {
        fmt.Println("getting terminal (to exit type exit):", emp_name)
        employes_container_id := FindContainerID(client, emp_name)
        if employes_container_id == "" {
                panic("employee's container not found")
        }
        exec, err := client.CreateExec(docker.CreateExecOptions{
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
        if err = dockerpty.StartExec(client, exec); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        fmt.Println("hope you enjoyed your stay dont forget to save your work as docker image")
}

func SaveImage(client *docker.Client, emp_name string) {
        fmt.Println("saving your image with name:", emp_name)
        employes_container_id := FindContainerID(client, emp_name)
        err := client.RemoveImage(emp_name)
        if err != nil {
                fmt.Println(err)
                //os.Exit(1)
        }

        img, err := client.CommitContainer(docker.CommitContainerOptions{
                Container:  employes_container_id,
                Author:     "docker-manager",
                Repository: emp_name,
                Tag:        "latest",
        })
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        fmt.Println("image saved:", img.ID)
}

func main() {
        endpoint := "unix:///var/run/docker.sock"
        client, err := docker.NewClient(endpoint)
        if err != nil {
                panic(err)
        }
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter employid: ")
        raw_input, _ := reader.ReadString('\n')
        emp_name := raw_input[:len(raw_input)-1]
        fmt.Println("trying to fetch your docker container")
        status := GetStatus(client, emp_name)
        fmt.Println("status container:", status)
        if status == "running" {
                GetTerminal(client, emp_name)
        } else if status == "stopped" {
                RunContainer(client, emp_name)
                GetTerminal(client, emp_name)
        } else if status == "exited" || status == "created" {
                RunContainer(client, emp_name)
                GetTerminal(client, emp_name)
        } else if status == "notfound" {
                CreateStartConatiner(client, emp_name)
                GetTerminal(client, emp_name)
        }
        SaveImage(client, emp_name)
}
