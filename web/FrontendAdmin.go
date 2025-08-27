package web

import (
	"GoodBuy/constants"
	"GoodBuy/db"
	"GoodBuy/security"
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

	currentRole := db.GetRolenameOfUserByName(currentUser)
	if currentRole != "Admin" {
		role_blocks := blocks(currentUser)

		data := map[string]any{
			"title": "Администрирование",
			"user":  currentUser,

			"current_year":     constants.CURRENT_YEAR(),
			"DEFAULT_PASSWORD": constants.DEFAULT_PASSWORD(),
		}

		t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/forbidden.html")
		err = t.Execute(w, data)
		if err != nil {
			println(err.Error())
		}
		return
	}

	postalert := ""

	if r.Method == http.MethodPost {
		if targetUser := r.PostFormValue("changeUser"); targetUser != "" {
			targetUser = strings.TrimPrefix(targetUser, "user")
			role, _ := strconv.Atoi(r.PostFormValue("role" + targetUser))
			password := r.PostFormValue("password" + targetUser)

			if strings.ToLower(targetUser) != "new" {
				user, _ := strconv.Atoi(targetUser)
				currentRole := db.GetRolenameOfUserById(uint8(user))

				if strings.Contains(strings.ToLower(currentRole), "admin") {
					postalert += "Нельзя изменять роль администратора через визуальный интерфейс.\n"
				} else if currentRole != db.GetRolenameByID(uint8(role)) {
					db.GrantRoleToUser(uint8(user), role)
				}

				if password != "" {
					db.ChangeUserPassword(uint8(user), security.Hash(password))
				}
			} else {
				username := r.PostFormValue("username" + targetUser)
				if password == "" {
					password = constants.DEFAULT_PASSWORD()
				}
				if !db.IsUserExist(username) {
					db.RegisterUser(uint8(role), username, security.Hash(password))
				}
			}
		}

		if targetUser := r.PostFormValue("deleteUser"); targetUser != "" {
			targetUser = strings.TrimPrefix(targetUser, "user")
			userID, err := strconv.Atoi(targetUser)
			if err != nil {
				postalert += "Ошибка в форме при вычислении ID пользователя: " + err.Error() + ".\n"
			} else if !db.IsUserExistByID(uint8(userID)) {
				postalert += "Пользователя с ID " + targetUser + " не существует.\n"
			} else if strings.Contains(strings.ToLower(db.GetRolenameOfUserById(uint8(userID))), "admin") {
				postalert += "Пользователь с ID " + targetUser + " является админом, чтобы отключить или полностью удалить админа из системы, сделайте соответствующий sql запрос.\n"
			} else {
				db.DisableUserByID(uint8(userID))
				postalert += "Пользователь " + targetUser + " отключен, чтобы полностью удалить его из системы, сделайте соответствующий sql запрос.\n"
			}
		}

		if targetProductCategory := r.PostFormValue("changeProductCategory"); targetProductCategory != "" {
			targetProductCategory = strings.TrimPrefix(targetProductCategory, "product_category")
			name := r.PostFormValue("pcname" + targetProductCategory)
			description := r.PostFormValue("pcdescription" + targetProductCategory)

			if strings.ToLower(targetProductCategory) != "new" {
				productID, err := strconv.Atoi(targetProductCategory)
				if err == nil {
					db.UpdateProductCategory(uint8(productID), name, description)
				} else {
					println("Something on getting product ID", err.Error())
				}
			} else {
				db.AddProductCategory(name, description)
			}
		}

		if r.PostFormValue("add_market") != "" {
			fee, _ := strconv.ParseFloat(r.PostFormValue("fee"), 64)
			db.AddMarket(r.PostFormValue("market"), r.PostFormValue("date_start"), r.PostFormValue("date_end"), fee)
		}
	}

	role_blocks := blocks(currentUser)

	data := map[string]any{
		"title":              "Администрирование",
		"user":               currentUser,
		"users":              db.GetUsers(),
		"roles":              db.GetRoles(),
		"product_categories": db.GetCategories(),
		"alertmessage":       postalert,

		"current_year":     constants.CURRENT_YEAR(),
		"DEFAULT_PASSWORD": constants.DEFAULT_PASSWORD(),
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendAdmin_main.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
