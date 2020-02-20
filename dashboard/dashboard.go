package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"text/template"
	"time"
)

type Pod struct {
	Name        string
	Ready       string
	Status      string
	Restarts    string
	Age         string
	Port        string
	Description string
}
type PodData struct {
	Podlist []Pod
}

var PodMaster PodData

func main() {
	go GrabPods()
	StartHTMLServer("4000")
}

func Pods(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./webpages/pods.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method == "POST" {
		inputname := r.FormValue("name")
		fmt.Println(inputname)
		inputimage := r.FormValue("image")
		fmt.Println(inputimage)
		inputport := r.FormValue("port")
		fmt.Println(inputport)
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte("kubectl run " + inputname + " --image=" + inputimage + " --port=" + inputport + "\n"))
		conn.Close()
		net.Dial("tcp", ":8080")
		conn.Write([]byte("kubectl expose deployment " + inputname + " --type=NodePort --name=" + inputname))
		conn.Close()
	}

	t.Execute(w, PodMaster)
}

func Nodes(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./webpages/nodes.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method == "POST" {
		inputname := r.FormValue("name")
		fmt.Println(inputname)
		inputimage := r.FormValue("image")
		fmt.Println(inputimage)
		inputport := r.FormValue("port")
		fmt.Println(inputport)
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte("kubectl run " + inputname + " --image=" + inputimage + " --port=" + inputport + "\n"))
		conn.Close()
		conn, _ = net.Dial("tcp", ":8080")
		conn.Write([]byte("kubectl expose deployment " + inputname + " --type=NodePort --name=" + inputname))
		conn.Close()
	}

	t.Execute(w, PodMaster)
}

func Deployments(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./webpages/deployments.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method == "POST" {
		inputname := r.FormValue("name")
		fmt.Println(inputname)
		inputimage := r.FormValue("image")
		fmt.Println(inputimage)
		inputport := r.FormValue("port")
		fmt.Println(inputport)
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte("kubectl run " + inputname + " --image=" + inputimage + " --port=" + inputport + "\n"))
		conn.Close()
		net.Dial("tcp", ":8080")
		conn.Write([]byte("kubectl expose deployment " + inputname + " --type=NodePort --name=" + inputname))
		conn.Close()
	}

	t.Execute(w, PodMaster)
}

func Services(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./webpages/services.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method == "POST" {
		inputname := r.FormValue("name")
		fmt.Println(inputname)
		inputimage := r.FormValue("image")
		fmt.Println(inputimage)
		inputport := r.FormValue("port")
		fmt.Println(inputport)
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte("kubectl run " + inputname + " --image=" + inputimage + " --port=" + inputport + "\n"))
		conn.Close()
		net.Dial("tcp", ":8080")
		conn.Write([]byte("kubectl expose deployment " + inputname + " --type=NodePort --name=" + inputname))
		conn.Close()
	}

	t.Execute(w, PodMaster)
}

//StartHTMLServer begins the hosting process for the
//webserver
func StartHTMLServer(port string) {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./webpages"))))
	http.HandleFunc("/", Pods)
	http.HandleFunc("/deployments", Deployments)
	http.HandleFunc("/services", Services)
	http.HandleFunc("/nodes", Nodes)
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
