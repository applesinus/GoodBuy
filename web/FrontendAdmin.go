package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func admin(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value

	currentRole := db.GetRoleOfUser(currentUser)
	if currentRole != "Admin" {
		role_blocks := blocks(currentUser)

		data := map[string]interface{}{
			"title": "Администрирование",
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
		if targetUser := r.PostFormValue("changeUser"); targetUser != "" {
			targetUser = strings.TrimPrefix(targetUser, "user")
			role, _ := strconv.Atoi(r.PostFormValue("role" + targetUser))
			user, _ := strconv.Atoi(targetUser)
			db.GrantRoleToUser(db.GetUsernameById(uint8(user)), role)
		}

		/*if targetProductCategory := r.PostFormValue("changeProductCategory"); targetProductCategory != "" {
			targetProductCategory = strings.TrimPrefix(targetProductCategory, "product_category")
			name, _ := strconv.Atoi(r.PostFormValue("product_category" + targetProductCategory))
			description, _ := strconv.Atoi(r.PostFormValue("description" + targetProductCategory))
			user, _ := strconv.Atoi(targetProductCategory)
			db.ChangeProductCategoryDescription(db.GetUsernameById(uint8(user)), role)
		}*/

		if r.PostFormValue("add_market") != "" {
			fee, _ := strconv.ParseFloat(r.PostFormValue("fee"), 64)
			db.AddMarket(r.PostFormValue("market"), r.PostFormValue("date_start"), r.PostFormValue("date_end"), fee)
		}
	}

	role_blocks := blocks(currentUser)

	data := map[string]interface{}{
		"title":              "Администрирование",
		"user":               currentUser,
		"users":              db.GetUsers(),
		"roles":              db.GetRoles(),
		"product_categories": db.GetCategories(),
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendAdmin_main.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
