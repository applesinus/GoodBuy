package web

import (
	"html/template"
	"net/http"
)

var reURL string

func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/login.html")
	t.Execute(w, nil)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/redirect.html")
	t.Execute(w, reURL)
}

func InitFront() {
	reURL = "/login"
	http.HandleFunc("/", redirect)
	http.HandleFunc("/login", login)
}
