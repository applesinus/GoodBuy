package web

import (
	"GoodBuy/db"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
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

	currentRole := db.GetRoleOfUser(currentUser)
	if currentRole != "Admin" && currentRole != "Salesman" {
		role_blocks := blocks(currentUser)

		data := map[string]interface{}{
			"title": "Товары",
			"user":  currentUser,
		}

		t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/forbidden.html")
		err = t.Execute(w, data)
		if err != nil {
			println(err.Error())
		}
		return
	}

	if r.Method == http.MethodPost {
		if r.PostFormValue("general") == "new" {
			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/products/new")
			return
		}
		if r.PostFormValue("change") != "" {
			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/products/edit?id="+r.PostFormValue("change"))
			return
		}
	}

	role_blocks := blocks(currentUser)

	pd := db.NewProduct()
	pds := db.GetProducts(pd, pd)

	data := map[string]interface{}{
		"title": "Товары",
		"user":  currentUser,
		"pds":   pds,
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendProducts_main.html")
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

	currentRole := db.GetRoleOfUser(currentUser)
	if currentRole != "Admin" && currentRole != "Salesman" {
		role_blocks := blocks(currentUser)

		data := map[string]interface{}{
			"title": "Добавить товар",
			"user":  currentUser,
		}

		t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/forbidden.html")
		err = t.Execute(w, data)
		if err != nil {
			println(err.Error())
		}
		return
	}

	if r.Method == http.MethodPost {
		new_product := db.NewProduct()
		new_product.Name = r.PostFormValue("Name")
		fmt.Sscan(r.PostFormValue("Amount"), &(new_product.Amount))
		fmt.Sscan(r.PostFormValue("Self_cost"), &(new_product.Self_cost))
		fmt.Sscan(r.PostFormValue("Default_cost"), &(new_product.Default_cost))
		fmt.Sscan(r.PostFormValue("Category"), &(new_product.Category))
		new_product.Category_name = db.GetCategoryNameById(new_product.Category)
		db.AddProduct(new_product)
	}

	categories := db.GetCategories()
	role_blocks := blocks(currentUser)

	data := map[string]interface{}{
		"title":      "Добавить товар",
		"user":       currentUser,
		"categories": categories,
		"alert":      "",
	}
	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendProducts_new.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}

func products_edit(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value

	currentRole := db.GetRoleOfUser(currentUser)
	if currentRole != "Admin" && currentRole != "Salesman" {
		role_blocks := blocks(currentUser)

		data := map[string]interface{}{
			"title": "Изменить товар",
			"user":  currentUser,
		}

		t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/forbidden.html")
		err = t.Execute(w, data)
		if err != nil {
			println(err.Error())
		}
		return
	}

	if r.Method == http.MethodPost {
		new_product := db.NewProduct()
		fmt.Sscan(r.PostFormValue("Id"), &(new_product.Id))
		new_product.Name = r.PostFormValue("Name")
		fmt.Sscan(r.PostFormValue("Amount"), &(new_product.Amount))
		fmt.Sscan(r.PostFormValue("Self_cost"), &(new_product.Self_cost))
		fmt.Sscan(r.PostFormValue("Default_cost"), &(new_product.Default_cost))
		fmt.Sscan(r.PostFormValue("Category"), &(new_product.Category))
		new_product.Category_name = db.GetCategoryNameById(new_product.Category)

		db.EditProduct(new_product.Id, new_product)

		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/products")
		return
	}

	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		println("Error parsing edit product id parameter from address.", err.Error())
		return
	}
	id, err := strconv.Atoi(queryParams.Get("id"))
	if err != nil {
		println("Error parsing edit product id parameter to int.", err.Error())
		return
	}
	product := db.GetProductByID(id)

	categories := db.GetCategories()
	role_blocks := blocks(currentUser)

	data := map[string]interface{}{
		"title":        "Изменить товар",
		"user":         currentUser,
		"categories":   categories,
		"Id":           product.Id,
		"Name":         product.Name,
		"Amount":       product.Amount,
		"Self_cost":    product.Self_cost,
		"Default_cost": product.Default_cost,
		"Category":     product.Category,
		"alert":        "",
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendProducts_edit.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
