package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"time"
)

const TIMETOSLEEP = 10 * time.Second

//Power is a control bool to be accessed to shut down the
//clientserver
var backendServers map[string]string = make(map[string]string)
var shutdownchan chan string = make(chan string)
var NodesList []Node

type Severs struct {
	ServerNames []string `json:"ServerNames"`
	NodePorts   []string `json:"NodePorts"`
}

type Node struct {
	Name    string
	Status  string
	Roles   string
	Age     string
	Version string
}

func main() {
	fmt.Println("Software Defined Network Terminal")
	go StartReverseProxy("80")
	go GetNodes()
	go GrabServers()
	<-shutdownchan
	fmt.Println("Shuting Down...")
	time.Sleep(TIMETOSLEEP)
}

//StartReverseProxy begins the hosting process for the
//client to server application
func StartReverseProxy(port string) {
	fmt.Println("Launching Software Defined Network...")

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		shutdownchan <- "Could Not Listen on Port"
		return
	}

	fmt.Println("Online - Now Listening On Port: " + port)

	ConnSignal := make(chan string)

	for {

		go Session(ln, ConnSignal, port)
		<-ConnSignal

	}

}

//Session creates a new seesion listening on a port. This
//session handles all interactions with the connected
//client
func Session(ln net.Listener, ConnSignal chan string, port string) {
	conn, _ := ln.Accept()
	defer conn.Close()
	ConnSignal <- "New Connection\n"

	//Checking for server to handle the connecting client
	buf := make([]byte, 1024)
	conn.Read(buf)
	var serverConn net.Conn = nil
	var err error
	for {
		for k, v := range backendServers {
			if strings.Contains(string(buf), k) {
				serverConn, err = net.Dial("tcp", v)
				if err != nil {
					//Server Could not be dialed, remove it from backendservers here
					return
				} else {
					defer serverConn.Close()
					serverConn.Write(buf)
					break
				}
			}
		}
		if serverConn != nil {
			break
		}

	}

	shutdownSession := make(chan string)
	go SessionListenerWriter(serverConn, conn, shutdownSession)
	go SessionListenerWriter(conn, serverConn, shutdownSession)
	<-shutdownSession
}

//SessionListenerWriter listens for connections noise and sends it to the writer
func SessionListenerWriter(Conn1 net.Conn, Conn2 net.Conn, shutdown chan string) {
	for {
		buf := make([]byte, 1024)
		Conn1.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, err := Conn1.Read(buf)
		if err != nil {
			fmt.Println(err)
			Conn1.Write([]byte("Timeout Error, No Signal. Disconnecting. \n"))
			break
		}
		Conn2.Write(buf)
	}
	shutdown <- "Session Closed"
}

//GrabServers allows user to add servers to list
func GrabServers() {
	for {
		openFile, _ := ioutil.ReadFile("serverlist.json")

		serverMain := Severs{}

		_ = json.Unmarshal(openFile, &serverMain)

		backendServers = make(map[string]string)
		for k, v := range serverMain.ServerNames {
			backendServers[v] = serverMain.NodePorts[k]
		}

		time.Sleep(TIMETOSLEEP)
	}

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
				TempNodesList = append(NodesList, NewNode)
			}

		}

		NodesList = TempNodesList

		byteslice, _ := json.MarshalIndent(NodesList, "", "	")

		ioutil.WriteFile("nodes.json", byteslice, 7777)

		time.Sleep(TIMETOSLEEP)
	}
}
