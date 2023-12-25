package main

import (
	"GoodBuy/db"
	"GoodBuy/web"
	"fmt"
	"strings"
	"time"

	"net/http"
)

const (
	isRunning = iota
	isStopped
	isErrored
)

func main() {

	db.StartDB()

	go web.InitBack()
	go web.InitFront()

	var server *http.Server

	serverStatus := isStopped

	command := "start"
	for strings.ToLower(command) != "exit" {
		switch strings.ToLower(command) {

		case "start":
			if serverStatus != isRunning {
				server = &http.Server{
					Addr: "localhost:8111",
				}
				serverIsRunning := make(chan bool)
				go func(serverIsRunning chan bool) {
					println("\nSERVER IS RUNNING!")
					serverIsRunning <- true
					err := server.ListenAndServe()
					if err != nil {
						println("\nTHE SERVER HAS BEEN STOPPED!")
						println(err.Error())
					}
				}(serverIsRunning)
				<-serverIsRunning
				close(serverIsRunning)
				serverStatus = isRunning
			}

		case "stop":
			server.Shutdown(nil)
			serverStatus = isStopped
			println("SERVER CLOSING IS DONE")

		case "exit":
			println("something's wrong, main.go case 'exit' (~line 60)")

		default:
			println("UNKNOWN COMMAND")

		}

		println("\n====================\nServer Status: ", serverStatus, "\nEnter the server command (start, stop, exit).")
		fmt.Scan(&command)
	}

	db.DropDB()

	time.Sleep(3 * time.Second)
}
