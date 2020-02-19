package main

import (
	"net"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {

	commandFile, _ := os.OpenFile("../commandlist", os.O_RDWR|os.O_CREATE, 7777)
	defer commandFile.Close()

	go commandListener(commandFile)

	readCommands(commandFile)
}

func readCommands(commandFile *os.File){
	for{
		AllBytes, _ := ioutil.ReadAll(commandFile)
		//BUG SPLAT: Possible command loss
		commandFile.Truncate(0)
		Lines := strings.Split(string(AllBytes), "\n")

		for _, v := range Lines {

			command := strings.Split(string(v), " ")
			out, _ := exec.Command(command[0], command[1:]...).Output()
			fmt.Print(string(out))
		}
	}
}
func commandListener(commandFile *os.File){
	//Pipe to control connection flow
	conPipe := make(chan string)

	ls, err := net.Listen("tcp", ":8080")
	defer ls.Close()
	if err != nil{
		fmt.Printf("Listen Error: %s", err)
	}
	
	//Cycle for loop as connections are made
	for{
		go commandConnection(ls, conPipe, commandFile)
		<-conPipe
	}
}
func commandConnection(ls net.Listener, conPipe chan string, commandFile *os.File){

	con, err := ls.Accept()
	defer con.Close()
	conPipe <- "Connection made"
	if err != nil {
		fmt.Printf("Accept error: %s", err)
		return
	}

	buf := make([]byte, 1024)
	con.Read(buf)
	
	commandFile.Write(buf)
}
 
