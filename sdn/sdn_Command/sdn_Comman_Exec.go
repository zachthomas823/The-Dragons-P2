package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {

	commandfile, _ := os.OpenFile("../commandlist", os.O_RDWR|os.O_CREATE, 7777)
	AllBytes, _ := ioutil.ReadAll(commandfile)
	commandfile.Truncate(0)
	Lines := strings.Split(string(AllBytes), "\n")

	for _, v := range Lines {

		command := strings.Split(string(v), " ")
		out, _ := exec.Command(command[0], command[1:]...).Output()
		fmt.Print(string(out))
	}
}
