package main

import (
	"github.com/NehemiahG7/project-1/Logger/Config"
	"fmt"
	"net"
	"log"
	"os"
)

func main(){
	file := loadFile(config.LogName)
	logger := log.New(file, config.LoggerPort, log.Flags())
	defer file.Close()

	ln, err := net.Listen("tcp", ":" + config.LoggerPort)
	if err != nil{
		logger.Fatalf("Listener failed: %s\n", err)
	}
	defer ln.Close()

	fmt.Printf("Logger listening on port :%s\n",config.LoggerPort)

	connect := make(chan string)

	for{
		go handleLog(ln, logger, connect)
		logger.Printf("Logger: connection from %s\n", <-connect)
	}
}
func handleLog(ln net.Listener, logger *log.Logger, connect chan string){
	for{
		conn, err := ln.Accept()
		defer conn.Close()
		if err != nil{
			logger.Printf("Logger: connection error: %s\n", err)
			return
		}

		logger.SetPrefix(conn.RemoteAddr().String() + " - ")

		connect <- string(conn.RemoteAddr().String())
		buff := make([]byte, 1024)

		num, err := conn.Read(buff)
		if err != nil{
			logger.Printf("Read failed: %s\n", err)
		}
		trm := buff[:num]

		logger.Printf("%s\n", trm)

	}
}
func loadFile(name string) *os.File{
	file, err := os.OpenFile(name, os.O_RDWR| os.O_APPEND| os.O_CREATE, 0666)
	if err != nil{
		log.Fatalf("Cannot open file: %s\n", err)
	}
	return file
}