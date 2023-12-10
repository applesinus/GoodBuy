package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
)

var reURL string

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
	var err error
	if isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/")
		return
	}

	if r.Method == http.MethodPost {
		login := r.PostFormValue("login")
		password := r.PostFormValue("password")

		if db.Auth(login, password) {
			reURL = "/products"

			cookie := &http.Cookie{
				Name:   "currentUser",
				Value:  login,
				MaxAge: 0,
			}
			http.SetCookie(w, cookie)
			cookie = &http.Cookie{
				Name:   "currentPassword",
				Value:  password,
				MaxAge: 0,
			}
			http.SetCookie(w, cookie)

			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/")
			return
		}
	}

	data := map[string]string{
		"title": "Вход",
		"user":  "",
	}

	logged_blocks, role_blocks := blocks(isLoggedIn(w, r), "")

	t, _ := template.ParseFiles("web/template.html", "web/"+logged_blocks, "web/"+role_blocks, "web/login.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	var err error
	t, _ := template.ParseFiles("web/redirect.html")
	err = t.Execute(w, reURL)
	if err != nil {
		println(err.Error())
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}

	cookie := &http.Cookie{
		Name:   "currentUser",
		Value:  "",
		MaxAge: 0,
	}
	http.SetCookie(w, cookie)
	cookie = &http.Cookie{
		Name:   "currentPassword",
		Value:  "",
		MaxAge: 0,
	}
	http.SetCookie(w, cookie)

	t, _ := template.ParseFiles("web/redirect.html")
	t.Execute(w, "/login")
	if err != nil {
		println(err.Error())
	}
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	user, erruser := r.Cookie("currentUser")
	password, errpassword := r.Cookie("currentPassword")

	if erruser != nil {
		cookie := &http.Cookie{
			Name:   "currentUser",
			Value:  "",
			MaxAge: 0,
		}
		http.SetCookie(w, cookie)

		if errpassword != nil {
			cookie = &http.Cookie{
				Name:   "currentPassword",
				Value:  "",
				MaxAge: 0,
			}
			http.SetCookie(w, cookie)
		}
		return false
	} else {
		if errpassword != nil {
			cookie := &http.Cookie{
				Name:   "currentPassword",
				Value:  "",
				MaxAge: 0,
			}
			http.SetCookie(w, cookie)
			return false
		}
	}

	if db.Auth(user.Value, password.Value) {
		return true
	} else {
		return false
	}

}

func InitFront() {
	reURL = "/login"
	http.HandleFunc("/", redirect)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/login", login)
	http.HandleFunc("/products", products)
	http.HandleFunc("/products/new", products_new)

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
}
