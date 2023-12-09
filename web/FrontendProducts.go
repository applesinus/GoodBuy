package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
)

func products(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}

	if r.Method == http.MethodPost {
		if r.PostFormValue("general") == "new" {
			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/products/new")
			return
		}
	}

	logged_blocks, role_blocks := blocks(isLoggedIn, currentUser)

	pd := db.NewProduct()
	pds := db.GetProducts(pd, pd)

	data := map[string]interface{}{
		"title": "Продукты",
		"user":  currentUser,
		"pds":   pds,
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+logged_blocks, "web/"+role_blocks, "web/products.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}

func products_new(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	logged_blocks, role_blocks := blocks(isLoggedIn, currentUser)

	// TODO

	categories := db.GetCategories()

	data := map[string]interface{}{
		"title":      "Добавить новый продукт",
		"user":       currentUser,
		"categories": categories,
	}
	t, _ := template.ParseFiles("web/template.html", "web/"+logged_blocks, "web/"+role_blocks, "web/products_new.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
