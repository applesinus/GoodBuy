package web

import (
	//"GoodBuy/db"
	"GoodBuy/db"
	"html/template"
	"net/http"
)

func receipts(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value
	logged_blocks, role_blocks := blocks(isLoggedIn(w, r), currentUser)

	if r.Method == http.MethodPost {
		if r.PostFormValue("general") == "new_receipt" {
			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/receipts/new")
			return
		}
	}

	//TODO
	rec1 := db.NewReceipt()
	rec2 := db.NewReceipt()

	data := map[string]interface{}{
		"title":    "Продукты",
		"user":     currentUser,
		"receipts": []db.Receipt{rec1, rec2},
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+logged_blocks, "web/"+role_blocks, "web/FrontendReceipts_main.html")
	t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}

func reciepts_new(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value
	logged_blocks, role_blocks := blocks(isLoggedIn(w, r), currentUser)

	pd := db.NewProduct()
	pds := db.GetProducts(pd, pd)

	data := map[string]interface{}{
		"title": "Продукты",
		"user":  currentUser,
		"pds":   pds,
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+logged_blocks, "web/"+role_blocks, "web/FrontendReceipts_new.html")
	t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
