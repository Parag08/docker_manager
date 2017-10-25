package main

import (
   "testing"
   "github.com/fsouza/go-dockerclient"
)

func contains(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
}

func TestGetStatus(t *testing.T) {
    endpoint := "unix:///var/run/docker.sock"
   client, err := docker.NewClient(endpoint)
   if err != nil {
       t.Fatal(err)
   }
   emp_name := "test"
   status := GetStatus(client, emp_name)
   validStatus := []string{"running","stopped","exited","created","notfound"}
   if !contains(validStatus,status) {
       t.Error("expected one of the %v",validStatus)
   }
}

func TestCreateStartConatiner(t *testing.T) {
   endpoint := "unix:///var/run/docker.sock"
   client, err := docker.NewClient(endpoint)
   if err != nil {
       t.Fatal(err)
   }
   emp_name := "test"
   status := GetStatus(client, emp_name)
   validStatus := []string{"running","stopped","exited","created","notfound"}
   if !contains(validStatus,status) {
       t.Error("expected one of the %v",validStatus)
   }
   if status != "notfound" {
       t.Skip("we cannot do this test at this time because the container need to be in notfound status to do this test")
   } else {
       CreateStartConatiner(client, emp_name)
   }

}

func TestRunContainer(t *testing.T) {
   endpoint := "unix:///var/run/docker.sock"
   client, err := docker.NewClient(endpoint)
   if err != nil {
       t.Fatal(err)
   }
   emp_name := "test"
   status := GetStatus(client, emp_name)
   validStatus := []string{"running","stopped","exited","created","notfound"}
   if !contains(validStatus,status) {
       t.Error("expected one of the %v",validStatus)
   }
   if status != "exited" || status != "created" {
       t.Skip("we cannot do this test at this time because the container need to be in exited/created status to do this test")
   } else {
       CreateStartConatiner(client, emp_name)
   }
}


func TestGetTerminal(t *testing.T) {
   endpoint := "unix:///var/run/docker.sock"
   client, err := docker.NewClient(endpoint)
   if err != nil {
       t.Fatal(err)
   }
   emp_name := "test"
   status := GetStatus(client, emp_name)
   validStatus := []string{"running","stopped","exited","created","notfound"}
   if !contains(validStatus,status) {
       t.Error("expected one of the %v",validStatus)
   }
   if status != "running" {
       t.Error("we cannot do this test at this time because the container need to be in ruuning status to do this test")
   } else {
       GetTerminal(client, emp_name)
   }
}


func TestFindContainerID(t *testing.T) {
   endpoint := "unix:///var/run/docker.sock"
   client, err := docker.NewClient(endpoint)
   if err != nil {
       t.Fatal(err)
   }
   emp_name := "test"
   employes_container_id := FindContainerID(client, emp_name)
   t.Log("container ID:",employes_container_id)
}


func TestSaveImage(t *testing.T) {
   endpoint := "unix:///var/run/docker.sock"
   client, err := docker.NewClient(endpoint)
   if err != nil {
       t.Fatal(err)
   }
   emp_name := "test"
   SaveImage(client, emp_name)
}
