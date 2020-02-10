package main

import(
	"github.com/NehemiahG7/project-1/ServerManager/Config"
	"github.com/NehemiahG7/project-1/ServerManager/balancer"
	"github.com/NehemiahG7/project-1/ServerManager/logger"
	"github.com/NehemiahG7/project-1/ServerManager/server"
	"github.com/NehemiahG7/project-1/ServerManager/rProxy"
	//"os/exec"
	//"log"
	"fmt"
	"net"
	//"strings"
)

//Channel to pass messages for the logger
var logCh chan string = make(chan string)

var cmdCh chan string = make(chan string)

var power chan string = make(chan string)


func main(){
	//Start Logger and wait for response
	go logger.StartLogger(cmdCh)
	fmt.Println(<- cmdCh)

	//Connect to logger
	go handleLog()
	logCh <- "Manager Online"

	//Start the rest of the server package
	go rproxy.StartProxy(logCh)
	go balancer.StartBalancer(config.NumServs, logCh)
	<- logCh
	for i := 0; i < config.NumServs; i++{
		go server.StartServer(config.Container, i, logCh)
	}
	<- power
}




func handleLog(){
	for{
		//conn, err := net.Dial("tcp", "127.0.0.1:" + config.LoggerPort)
		conn, err := net.Dial("tcp", "localhost:" + config.LoggerPort)
		if err != nil{
			fmt.Printf("Logger not loaded %s\n", err)
			continue
		}
		conn.Write([]byte(<-logCh))
	}
	
}
