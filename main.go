package main

import (
	"GoodBuy/db"
	"GoodBuy/web"
	"fmt"
	"strings"
	"time"

	"net/http"
)

func main() {

	db.ConnectDB()
	println("Is this your first time of using the program on this PC?")
	var line string
	fmt.Scan(&line)
	if strings.ToLower(line) == "yes" {
		db.CreateDB()
	} else {
		db.StartDB()
	}

	go web.InitBack()
	go web.InitFront()

	var server *http.Server
	server = &http.Server{
		Addr: "localhost:8111",
	}

	go func() {
		println("\nSERVER IS RUNNING!")
		err := server.ListenAndServe()
		if err != nil {
			println("\nSERVER HAS BEEN STOPPED!")
		}
	}()

	time.Sleep(3 * time.Second)

	print("Enter the server command: ")
	var command string
	fmt.Scan(&command)
	if strings.ToLower(command) == "stop" {
		server.Close()
	}
	println("IT'S DONE")

	time.Sleep(3 * time.Second)
}
