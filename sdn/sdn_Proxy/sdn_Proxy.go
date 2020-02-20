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
var PublicIP string

func main() {
	fmt.Println("Software Defined Network Terminal")
	go GrabServers()
	go StartReverseProxy("4000")
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
				serverConn, err = net.Dial("tcp", ":"+v)
				if err != nil {
					conn.Write([]byte("Could not resolve: " + PublicIP + ":" + v))
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

	toclient := make(chan []byte)
	toserver := make(chan []byte)
	shutdownSession := make(chan string)
	go SessionWriter(conn, shutdownchan, toclient)
	go SessionWriter(serverConn, shutdownchan, toserver)
	go SessionListener(serverConn, shutdownSession, toclient)
	go SessionListener(conn, shutdownSession, toserver)
	<-shutdownSession
	time.Sleep(5 * time.Second)
}

//SessionListener listens for connections noise and sends it to the writer
func SessionListener(Conn net.Conn, shutdown chan string, writer chan []byte) {
	for {
		buf := make([]byte, 1024)
		Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, err := Conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}
		writer <- buf
	}
	shutdown <- "Session Closed"
}

func SessionWriter(Conn net.Conn, shutdown chan string, writer chan []byte) {

	for {
		buf := <-writer

		Conn.Write(buf)

	}
}

//GrabServers test
func GrabServers() {
	for {

		openFile, _ := ioutil.ReadFile("../serverlist.json")

		_ = json.Unmarshal(openFile, &backendServers)

		time.Sleep(TIMETOSLEEP)
	}
}
