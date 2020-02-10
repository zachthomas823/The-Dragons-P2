package config

import (
	"flag"
)

//LoggerPort holds the value for the pP flag
var LoggerPort string

//LogName holds the value for the lN flag
var LogName string

func init(){
	flag.StringVar(&LoggerPort, "lP", "9090", "Use lP to designate the port that the logger listens on")
	flag.StringVar(&LogName, "lN", "txt.log", "Use lN to designate the name of the log file")
	
	flag.Parse()
}