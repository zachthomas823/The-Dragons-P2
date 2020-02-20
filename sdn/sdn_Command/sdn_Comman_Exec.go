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
	go commandListener("8080")

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
			if v != 0 {
				temp = append(temp, v)
			}
		}
		AllBytes = temp

		Lines := strings.Split(string(AllBytes), "\n")

		for _, v := range Lines {
			command := strings.Split(string(v), " ")
			_, err := exec.Command(command[0], command[1:]...).Output()
			if err != nil {
				fmt.Println(err)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func commandListener(port string) {
	//Pipe to control connection flow
	commandFile, _ := os.OpenFile("../commandlist", os.O_RDWR|os.O_CREATE, 7777)
	conPipe := make(chan string)

	ls, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Printf("Listen Error: %s", err)
	}

	defer ls.Close()

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
