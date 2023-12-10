package web

import (
	"GoodBuy/db"
	"fmt"
	"html/template"
	"net/http"
)

func products(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value

	if r.Method == http.MethodPost {
		if r.PostFormValue("general") == "new" {
			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/products/new")
			return
		}
	}

	logged_blocks, role_blocks := blocks(isLoggedIn(w, r), currentUser)

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
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value

	if r.Method == http.MethodPost {
		new_product := db.NewProduct()
		new_product.Name = r.PostFormValue("Name")
		fmt.Sscan(r.PostFormValue("Amount"), &(new_product.Amount))
		fmt.Sscan(r.PostFormValue("Self_cost"), &(new_product.Self_cost))
		fmt.Sscan(r.PostFormValue("Default_cost"), &(new_product.Default_cost))
		fmt.Sscan(r.PostFormValue("Category"), &(new_product.Category))
		new_product.Category_name = db.GetCategotyNameById(new_product.Category)
		db.AddProduct(new_product)
	}

	categories := db.GetCategories()
	logged_blocks, role_blocks := blocks(isLoggedIn(w, r), currentUser)

	data := map[string]interface{}{
		"title":      "Добавить новый продукт",
		"user":       currentUser,
		"categories": categories,
		"alert":      "",
	}
	t, _ := template.ParseFiles("web/template.html", "web/"+logged_blocks, "web/"+role_blocks, "web/products_new.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
