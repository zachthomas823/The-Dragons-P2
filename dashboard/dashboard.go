package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"
)

type Pod struct {
	Name     string
	Ready    string
	Status   string
	Restarts string
	Age      string
	Port     string
}
type PodData struct {
	Podlist []Pod
}

var PodMaster PodData

func main() {
	go GrabPods()
	StartHTMLServer("8080")
}

//Handler is main handler for running the webpage.
//It passes the global variable Database for parsing
//index.html for pulling data.
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(PodMaster.Podlist)
	t, err := template.ParseFiles("./webpages/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, PodMaster)
}

//StartHTMLServer begins the hosting process for the
//webserver
func StartHTMLServer(port string) {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./webpages"))))
	http.HandleFunc("/", Handler)
	fmt.Println("Online - Now Listening On Port: " + port)

	http.ListenAndServe(":"+port, nil)
}

func GrabPods() {
	for {
		openFile, _ := ioutil.ReadFile("../sdn/pods.json")

		_ = json.Unmarshal(openFile, &PodMaster.Podlist)

		time.Sleep(10 * time.Second)
	}
}
