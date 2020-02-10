package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"time"
)

//Power is a control bool to be accessed to shut down the
//clientserver
var Power bool = true
var connectionsPerServer map[string]int = make(map[string]int)
var shutdownchan chan string
var logConn net.Conn
var port string

func main() {
	fmt.Println("Load Balancer Terminal")
	go GrabServers()
	<-shutdownchan
}

//StartLoadBalancer begins the hosting process for the
//load balancer which assumes all incomming traffic is
//for the same type of server and routes messages to the
//least used server
func StartLoadBalancer(port string) {
	fmt.Println("Launching Load Balancing server...")

	ln, _ := net.Listen("tcp", ":"+port)

	fmt.Println("Online - Now Listening On Port: " + port)
	fmt.Println()

	ConnSignal := make(chan string)

	for Power {

		go Session(ln, ConnSignal, port)
		<-ConnSignal

	}
	fmt.Println("Shut Down Signal Sent...Ending")
}

var shutDownSession chan string = make(chan string)

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

	//Matches conn to initial Server based on server conns
	//through the loadbalancer.
	var serverConn net.Conn = nil
	var intslice []int
	var err error

	for k := range connectionsPerServer {
		intslice = append(intslice, connectionsPerServer[k])
	}
	sort.Ints(intslice)
	valueToLookfor := intslice[0]

	for i := 0; i < 10; i++ {
		for k, v := range connectionsPerServer {
			if valueToLookfor == v {

				serverConn, err = net.Dial("tcp", k)

				if err != nil {
					fmt.Println(k)
					fmt.Println("Server Not Responding with Error, Removing from list")
					fmt.Println(err)
					temp := make(map[string]int)
					for key, value := range connectionsPerServer {
						if key != k {
							temp[key] = value
						}
					}
					connectionsPerServer = temp
					return
				}

				defer serverConn.Close()
				serverConn.Write(buf)
				connectionsPerServer[k]++
				break

			}
		}
		if serverConn != nil {
			break
		}
	}

	Shutdown := make(chan string)
	InboundMessages := make(chan []byte)
	OutboundMessages := make(chan []byte)
	go SessionWriter(conn, OutboundMessages, Shutdown)
	go SessionWriter(serverConn, InboundMessages, Shutdown)
	go SessionListener(serverConn, OutboundMessages, Shutdown)
	go SessionListener(conn, InboundMessages, Shutdown)
	<-Shutdown

	connectionsPerServer[serverConn.RemoteAddr().String()] = connectionsPerServer[serverConn.RemoteAddr().String()] - 1
}

//SessionListener listens for connections noise and sends it to the writer
func SessionListener(Conn1 net.Conn, messages chan []byte, shutdown chan string) {
	for {
		buf := make([]byte, 1024)
		Conn1.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, err := Conn1.Read(buf)
		if err != nil {
			fmt.Println(err)
			Conn1.Write([]byte("Timeout Error, No Signal. Disconnecting. \n"))
			break
		}

		if logConn != nil {
			var logfile []byte
			if strings.Contains(Conn1.RemoteAddr().String(), port) {
				logfile = []byte("In <" + Conn1.LocalAddr().String() + ">\n")
			} else {
				logfile = []byte("In <" + Conn1.RemoteAddr().String() + ">\n")
			}
			logConn.Write(append(logfile, buf...))
		}

		messages <- buf
	}
	shutdown <- "Close Session"
}

//SessionWriter listens for messages channel and sends it to the correct server
func SessionWriter(Conn1 net.Conn, messages chan []byte, shutdown chan string) {
	for {
		NewMessage := <-messages

		if logConn != nil {
			var logfile []byte
			if strings.Contains(Conn1.RemoteAddr().String(), port) {
				logfile = []byte("OUT <" + Conn1.LocalAddr().String() + ">\n")
			} else {
				logfile = []byte("OUT <" + Conn1.RemoteAddr().String() + ">\n")
			}
			logConn.Write(append(logfile, NewMessage...))
		}

		Conn1.Write(NewMessage)
	}
}

//GrabServers allows user to add servers to list
func GrabServers() {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(string(text), " \n")
		text = strings.TrimSpace(string(text))

		switch text {
		case "add":
			fmt.Println("Grab Servers By Entering in a full address such as Host:Port")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimRight(string(text), " \n")
			text = strings.TrimSpace(string(text))
			conn := text

			connectionsPerServer[conn] = 0
			fmt.Println("Added")

		case "Log":
			fmt.Println("Grab logging server By Entering in a full address such as Host:Port")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimRight(string(text), " \n")
			text = strings.TrimSpace(string(text))

			conn, _ := net.Dial("tcp", text)
			logConn = conn

		case "Launch":
			fmt.Println("Enter in a port")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimRight(string(text), " \n")
			text = strings.TrimSpace(string(text))

			port = text

			go StartLoadBalancer(text)

		case "Remove":
			fmt.Println("Remove Server By Entering in a full address such as Host:Port")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimRight(string(text), " \n")
			text = strings.TrimSpace(string(text))

			temp := make(map[string]int)
			for k, v := range connectionsPerServer {
				if k != text {
					temp[k] = v
				}
			}
			connectionsPerServer = temp
		}

	}

}
