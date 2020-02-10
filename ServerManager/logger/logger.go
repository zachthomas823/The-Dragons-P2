package logger

import (
	"log"
	"os/exec"
)

//StartLogger starts the logger for the package
func StartLogger(cmdCh chan string){
	cmd := exec.Command("./Logger/logger")
	cmd.Start()
	
	cmdCh <- "Logger started"
	
	err := cmd.Wait()
	if err != nil{
		log.Fatalf("Logger could not be started, closing program: %s\n", err)
	}
}