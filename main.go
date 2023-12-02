package main

import (
	"GoodBuy/db"
	"GoodBuy/web"

	"net/http"
)

func main() {

	go db.StartDB()

	go web.InitBack()
	go web.InitFront()

	http.ListenAndServe(":8100", nil)

}
