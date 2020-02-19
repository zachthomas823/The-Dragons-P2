package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const TIMETOSLEEP = 10 * time.Second

type Node struct {
	Name    string
	Status  string
	Roles   string
	Age     string
	Version string
}

type Pod struct {
	Name        string
	Ready       string
	Status      string
	Restarts    string
	Age         string
	Port        string
	Description string
}

type Service struct {
	Name       string
	Type       string
	ClusterIP  string
	ExternalIP string
	Port       string
	Age        string
}

func main() {
	go GetNodes()
	go GetPods()
	GetServices()
}

func GetNodes() {
	for {
		var NewNode Node
		var TempNodesList []Node

		output, _ := exec.Command("kubectl", "get", "nodes").Output()

		t := strings.Split(string(output), "\n")
		t = t[1:]

		for _, v := range t {
			z := strings.Split(v, " ")

			var temp []string

			for k2, v2 := range z {
				z[k2] = strings.TrimSpace(v2)
				if z[k2] != "" {
					temp = append(temp, z[k2])
				}
			}

			z = temp

			if len(z) != 0 {
				NewNode = Node{Name: z[0], Status: z[1], Roles: z[2], Age: z[3], Version: z[4]}
				TempNodesList = append(TempNodesList, NewNode)
			}

		}

		byteslice, _ := json.MarshalIndent(TempNodesList, "", "	")

		ioutil.WriteFile("../nodes.json", byteslice, 7777)

		time.Sleep(TIMETOSLEEP)
	}
}

func GetPods() {
	for {
		var NewPod Pod
		var TempPodList []Pod

		output, _ := exec.Command("kubectl", "get", "pods").Output()

		t := strings.Split(string(output), "\n")
		t = t[1:]

		for _, v := range t {
			z := strings.Split(v, " ")

			var temp []string

			for k2, v2 := range z {
				z[k2] = strings.TrimSpace(v2)
				if z[k2] != "" {
					temp = append(temp, z[k2])
				}
			}

			z = temp

			if len(z) != 0 {
				descrip, _ := exec.Command("kubectl", "describe", "pods", z[0]).Output()
				descripfile, _ := os.OpenFile("./pods", os.O_RDWR|os.O_CREATE, 7777)
				descripfile.Write(descrip)

				grepport, _ := exec.Command("grep", "Port:", "./pods").Output()
				grepportslice := strings.Split(string(grepport), "\n")
				grepportslice = strings.Split(string(grepportslice[0]), " ")

				port := grepportslice[len(grepportslice)-1]
				portslice := strings.Split(port, "/")

				NewPod = Pod{Name: z[0], Ready: z[1], Status: z[2], Restarts: z[3], Age: z[4], Port: portslice[0], Description: string(descrip)}
				TempPodList = append(TempPodList, NewPod)
			}

		}

		byteslice, _ := json.MarshalIndent(TempPodList, "", "	")

		ioutil.WriteFile("../pods.json", byteslice, 7777)

		time.Sleep(TIMETOSLEEP)
	}
}

func GetServices() {
	for {
		var NewService Service
		var TempNewServiceList []Service

		output, _ := exec.Command("kubectl", "get", "svc").Output()

		t := strings.Split(string(output), "\n")
		t = t[1:]

		for _, v := range t {
			z := strings.Split(v, " ")

			var temp []string

			for k2, v2 := range z {
				z[k2] = strings.TrimSpace(v2)
				if z[k2] != "" {
					temp = append(temp, z[k2])
				}
			}

			z = temp

			if len(z) != 0 {
				NewService = Service{Name: z[0], Type: z[1], ClusterIP: z[2], ExternalIP: z[3], Port: z[4], Age: z[5]}
				TempNewServiceList = append(TempNewServiceList, NewService)
			}

		}

		byteslice, _ := json.MarshalIndent(TempNewServiceList, "", "	")

		ioutil.WriteFile("../services.json", byteslice, 7777)

		time.Sleep(TIMETOSLEEP)
	}
}
