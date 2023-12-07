package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
)

var reURL string

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.PostFormValue("login")
		password := r.PostFormValue("password")

		if db.Auth(login, password) {
			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/success")
		}
	}

	t, _ := template.ParseFiles("web/login.html")
	t.Execute(w, nil)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/redirect.html")
	t.Execute(w, reURL)
}

func success(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func InitBack() {
	reURL = "/login"
	http.HandleFunc("/", redirect)
	http.HandleFunc("/login", login)
	http.HandleFunc("/success", success)
}
