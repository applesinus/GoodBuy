package web

import (
	"GoodBuy/db"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

func getIntArg(name string, def int, r *http.Request) int {
	if val := r.URL.Query().Get(name); val != "" {
		def, _ = strconv.Atoi(val)
	}
	return def
}

func buildURL(query url.Values, newArgs map[string]string, anchor string) string {
	newQuery := make(url.Values)
	for k, v := range query {
		newQuery[k] = v
	}
	for k, v := range newArgs {
		newQuery[k] = []string{v}
	}
	return newQuery.Encode() + "#" + anchor
}

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

	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		println("Error parsing edit product id parameter from address.", err.Error())
		return
	}

	if r.Method == http.MethodPost {
		if r.PostFormValue("total_days") != "" {
			http.Redirect(w, r,
				"/stats?"+buildURL(
					queryParams,
					map[string]string{
						"total_days": r.PostFormValue("total_days"),
					},
					"income",
				),
				http.StatusSeeOther,
			)
		} else if r.PostFormValue("mode") != "" {
			http.Redirect(w, r,
				"/stats?"+buildURL(
					queryParams,
					map[string]string{
						"mode_days": r.PostFormValue("mode_days"),
						"mode_n":    r.PostFormValue("mode_n"),
					},
					"mode",
				),
				http.StatusSeeOther,
			)
		} else if r.PostFormValue("profit") != "" {
			http.Redirect(w, r,
				"/stats?"+buildURL(
					queryParams,
					map[string]string{
						"profit_days": r.PostFormValue("profit_days"),
						"profit_n":    r.PostFormValue("profit_n"),
					},
					"profit",
				),
				http.StatusSeeOther,
			)
		} else if r.PostFormValue("mode_on_markets") != "" {
			http.Redirect(w, r,
				"/stats?"+buildURL(
					queryParams,
					map[string]string{
						"mode_on_markets_n": r.PostFormValue("mode_on_markets_n"),
					},
					"mode_on_markets",
				),
				http.StatusSeeOther,
			)
		}
	}

	total_days := getIntArg("total_days", 30, r)
	mode_days := getIntArg("mode_days", 30, r)
	mode_n := getIntArg("mode_n", 3, r)
	mode := db.GetMode(mode_days, mode_n)
	profit_days := getIntArg("profit_days", 30, r)
	profit_n := getIntArg("profit_n", 3, r)
	profit := db.GetMostProfitable(profit_days, profit_n)
	mode_on_markets_n := getIntArg("mode_on_markets_n", 3, r)
	mode_on_markets := db.GetModeOnMarkets(mode_on_markets_n)

	users := db.GetUsers()

	data := map[string]interface{}{
		"title":             "Статистика",
		"user":              currentUser,
		"users":             users,
		"roles":             db.GetRoles(),
		"income_by_days":    db.GetIncomePastNDays(total_days),
		"total_days":        total_days,
		"mode_days":         mode_days,
		"mode_n":            mode_n,
		"mode":              mode,
		"profit_days":       profit_days,
		"profit_n":          profit_n,
		"profit":            profit,
		"mode_on_markets_n": mode_on_markets_n,
		"mode_on_markets":   mode_on_markets,
	}

	t, _ := template.ParseFiles("web/template.html", "web/"+role_blocks, "web/FrontendStats_main.html")
	err = t.Execute(w, data)
	if err != nil {
		println(err.Error())
	}
}
