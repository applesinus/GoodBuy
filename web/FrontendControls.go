package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
)

var reURL string

func blocks(user string) string {

	var role_blocks string
	switch db.GetRoleOfUser(user) {
	case "Admin":
		role_blocks = "blocks_admin.html"
	case "Salesman":
		role_blocks = "blocks_seller.html"
	case "Analyst":
		role_blocks = "blocks_analyst.html"
	case "error":
		role_blocks = "blocks_notuser.html"
	}

	return role_blocks
}

func login(w http.ResponseWriter, r *http.Request) {
	var err error
	if isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/products")
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

	role_blocks := blocks("")

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/login.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	var err error
	if isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/products")
		return
	}

	if r.Method == http.MethodPost {
		login := r.PostFormValue("login")
		password := r.PostFormValue("password")

		if !db.Auth(login, password) {
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

			db.RegisterUser(login, password)

			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/")
			return
		}
	}

	data := map[string]string{
		"title": "Регистрация",
	}

	role_blocks := blocks("")

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/register.html")
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
	http.HandleFunc("/register", register)

	http.HandleFunc("/products", products)
	http.HandleFunc("/products/new", products_new)
	http.HandleFunc("/products/edit", products_edit)

	http.HandleFunc("/receipts", receipts)
	http.HandleFunc("/receipts/new", reciepts_new)

	http.HandleFunc("/admin", admin)

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
}
