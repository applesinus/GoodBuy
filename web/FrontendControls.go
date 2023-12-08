package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
)

var reURL string
var currentUser = ""
var isLoggedIn = false

func blocks(isLogged bool, user string) (string, string) {

	var logged_blocks string
	if isLogged {
		logged_blocks = "user_blocks.html"
	} else {
		logged_blocks = "notuser_blocks.html"
	}

	var role_blocks string
	if user == "Admin" {
		role_blocks = "admin_blocks.html"
	} else {
		role_blocks = "seller_blocks.html"
	}

	return logged_blocks, role_blocks
}

func login(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/")
		return
	}

	if r.Method == http.MethodPost {
		login := r.PostFormValue("login")
		password := r.PostFormValue("password")

		if db.Auth(login, password) {
			reURL = "/products"
			isLoggedIn = true
			currentUser = login
			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/")
			return
		}
	}

	data := map[string]string{
		"title": "Login",
		"user":  currentUser,
	}

	logged_blocks, role_blocks := blocks(isLoggedIn, currentUser)

	t, _ := template.ParseFiles("web/template.html", "web/"+logged_blocks, "web/"+role_blocks, "web/login.html")
	t.Execute(w, data)
}

func products(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/redirect.html")
	t.Execute(w, reURL)
}

func InitFront() {
	reURL = "/login"
	http.HandleFunc("/", redirect)
	http.HandleFunc("/login", login)
	http.HandleFunc("/products", products)

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
}
