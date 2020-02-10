package rproxy

import (
	"strings"
	"fmt"
	"os/exec"
)

//StartProxy starts a reverse proxy pointing to the load balancer
func StartProxy(logCh chan string){
	//Build command string
	str := "docker run --rm -p 8081:8081 --name proxy --network my-net proxy"

	//Split command string and creat cmd
	args := strings.Split(str, " ")
	cmd := exec.Command(args[0], args[1:]...)
	//cmd.Start()
	
	fmt.Println("Proxy started")
	// err := cmd.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil{
		logCh <- "Manager: proxy failed" + err.Error()
	} else {
		logCh <- "Manager: " + string(out) 
	}
}