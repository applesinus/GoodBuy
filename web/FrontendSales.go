package web

import (
	//"GoodBuy/db"
	"html/template"
	"net/http"
)

func sales(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value
	logged_blocks, role_blocks := blocks(isLoggedIn(w, r), currentUser)

	//TODO

	data := map[string]interface{}{
		"title": "Продукты",
		"user":  currentUser,
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+logged_blocks, "web/"+role_blocks, "web/sales.html")
	t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
