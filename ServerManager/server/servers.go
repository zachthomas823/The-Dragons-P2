package server

import (
	"fmt"
	"os/exec"
	"strings"
	"strconv"
)

//StartServer starts a server of the given container
func StartServer(cont string, num int, logCh chan string){

	s := "docker run --rm -d --name serv" + strconv.Itoa(num) + " --network my-net "+ cont
	logCh <- "Manager: " + s
	args := strings.Split(s, " ")
	
	//Create command with args[0] as command and rest as args
	cmd := exec.Command(args[0], args[1:]...)
	// cmd.Start()
	
	fmt.Println("Server started")
	// err := cmd.Wait()
	str, err := cmd.CombinedOutput()
	if err != nil{
		logCh <- "server could not be started, closing program: %s\n" + err.Error()
	} else {
		logCh <- "Manager: " + string(str) 
	}
}