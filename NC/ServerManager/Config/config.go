package config

import (
	"flag"
)

//StartCommands is a slice of strings that refer to all the commands needed to launch the server
var StartCommands []string
//LoggerPort holds the port value that points to the logger program
var LoggerPort string = "9090"
//NumServs hold the value for the n flag
var NumServs int
//Container holds the value for the c flag
var Container string

func init(){
	flag.IntVar(&NumServs, "n", 1, "Use n to set the number of server containers to run on startup")
	flag.StringVar(&Container, "c", "tictactoe", "Use the c flag to set the container to build servers from")
	
	flag.Parse()
}
//GetCommands returns the StartCommands slice
func GetCommands()[]string{
	return StartCommands
}