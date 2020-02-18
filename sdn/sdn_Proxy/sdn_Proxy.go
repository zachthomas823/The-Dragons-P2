package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

const TIMETOSLEEP = 10 * time.Second

//Power is a control bool to be accessed to shut down the
//clientserver
var backendServers map[string]string = make(map[string]string)
var shutdownchan chan string = make(chan string)

func main() {
	fmt.Println("Software Defined Network Terminal")

	go GrabServers()
	go StartReverseProxy("8080")
	<-shutdownchan
	fmt.Println("Shuting Down...")
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
	ConnSignal <- "New Connection \n"

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

//GrabServers test
func GrabServers() {

	openFile, _ := ioutil.ReadFile("./serverlist.json")

	_ = json.Unmarshal(openFile, &backendServers)

	time.Sleep(TIMETOSLEEP)

}
