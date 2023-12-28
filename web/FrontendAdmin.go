package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
	"strconv"
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
		if r.PostFormValue("change") != "" {
			role, _ := strconv.Atoi(r.PostFormValue("role" + r.PostFormValue("change")))
			user, _ := strconv.Atoi(r.PostFormValue("change"))
			db.GrantRoleToUser(db.GetUsernameById(uint8(user)), role)
		}
	}

	role_blocks := blocks(currentUser)

	users := db.GetUsers()

	data := map[string]interface{}{
		"title": "Администрирование",
		"user":  currentUser,
		"users": users,
		"roles": db.GetRoles(),
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendAdmin_main.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
