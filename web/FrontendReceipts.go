package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type date_str struct {
	Y int
	M time.Month
	D int
}

func receipts(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value

	currentRole := db.GetRolenameOfUserByName(currentUser)
	if currentRole != "Admin" && currentRole != "Salesman" {
		role_blocks := blocks(currentUser)

		data := map[string]interface{}{
			"title": "Продажи",
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

	if r.Method == http.MethodPost {
		if r.PostFormValue("general") == "new_receipt" {
			t, _ := template.ParseFiles("web/redirect.html")
			t.Execute(w, "/receipts/new")
			return
		}
	}

	rcps := db.GetAllReceipts()

	for i, rcp := range rcps {
		rcp.Pos_len++
		rcps[i] = rcp
	}

	data := map[string]interface{}{
		"title":    "Продажи",
		"user":     currentUser,
		"receipts": rcps,
	}

	t, err := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendReceipts_main.html")
	t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}

func receipts_new(w http.ResponseWriter, r *http.Request) {
	var err error
	if !isLoggedIn(w, r) {
		t, _ := template.ParseFiles("web/redirect.html")
		t.Execute(w, "/login")
		return
	}
	user, _ := r.Cookie("currentUser")
	currentUser := user.Value

	currentRole := db.GetRolenameOfUserByName(currentUser)
	if currentRole != "Admin" && currentRole != "Salesman" {
		role_blocks := blocks(currentUser)

		data := map[string]interface{}{
			"title": "Оформление продажи",
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
		i := 1
		receipt := db.NewReceipt()
		receipt.Date = time.Now().Format("2006-01-02")

		for {
			if r.PostFormValue("product"+strconv.Itoa(i)) != "" {
				position := db.NewPosition()
				position.Product = r.PostFormValue("product" + strconv.Itoa(i))
				cost, _ := strconv.ParseFloat(r.PostFormValue("cost"+strconv.Itoa(i)), 64)
				position.Cost = float32(cost)
				count, _ := strconv.ParseUint(r.PostFormValue("count"+strconv.Itoa(i)), 10, 8)
				position.Count = uint8(count)
				position.Status = "OK"
				receipt.Positions = append(receipt.Positions, position)
				receipt.Pos_len++
			} else {
				break
			}
			i++
		}
		db.AddNewReceipt(receipt)
	}

	role_blocks := blocks(currentUser)

	pd := db.NewProduct()
	pds := db.GetProducts(pd, pd)
	var date date_str
	date.Y, date.M, date.D = time.Now().Date()

	data := map[string]interface{}{
		"title": "Оформление продажи",
		"user":  currentUser,
		"Date":  date,
		"pds":   pds,
	}

	t, err := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendReceipts_new.html")
	t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
