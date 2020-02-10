package balancer

import (
	"fmt"
	"os/exec"
	"strings"
)

//StartBalancer starts a balancer container alternating between the number of servrs given
func StartBalancer(num int, logCh chan string){
	//create command string
	var b strings.Builder
	fmt.Fprintf(&b, "docker run --rm --name balancer --network my-net balancer -addr=")
	for i := 0; i < num; i++{
		if i == 0{
			fmt.Fprintf(&b, "serv%d", i)
		}else{
			fmt.Fprintf(&b, ",serv%d", i)
		}
	}
	logCh <- "Bal command: " + b.String()
	//split command string into command and args
	args := strings.Split(b.String(), " ")

	//Create command with args[0] as command and rest as args
	cmd := exec.Command(args[0], args[1:]...)
	// cmd.Start()
	
	fmt.Println("Balancer started")
	// err := cmd.Wait()
	str, err := cmd.CombinedOutput()
	if err != nil{
		logCh <- "Proxy could not be started, closing program: %s\n" + err.Error()
	} else {
		logCh <- "Manager: " + string(str) 
	}
}