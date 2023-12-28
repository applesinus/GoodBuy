package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
)

func stats(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value

	currentRole := db.GetRoleOfUser(currentUser)
	if currentRole != "Admin" && currentRole != "Analyst" {
		role_blocks := blocks(currentUser)

		data := map[string]interface{}{
			"title": "Статистика",
			"user":  currentUser,
		}

		t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/forbidden.html")
		err = t.Execute(w, data)
		if err != nil {
			println(err.Error())
		}
		return
	}

	role_blocks := blocks(currentUser)

	users := db.GetUsers()

	data := map[string]interface{}{
		"title":  "Статистика",
		"user":   currentUser,
		"users":  users,
		"roles":  db.GetRoles(),
		"income": db.GetAllIncome(),
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendStats_main.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
