package db

import "context"

type Mode struct {
	Name  string
	Sales int
}

type Profit struct {
	Name   string
	Profit float64
}

func GetIncomePastNDays(n int) int {
	sum := 0
	err := conn.QueryRow(context.Background(), "select * from goodbuy.get_income_past_N_days($1)", n).Scan(&sum)
	if err != nil {
		println("error getting total.", err.Error())
	}
	return sum
}

func GetMode(past_days int, n int) []Mode {
	var mode []Mode
	rows, err := conn.Query(context.Background(), "select * from goodbuy.get_top_N_products_by_sales($1, $2)", past_days, n)
	if err != nil {
		println("error getting mode.", err.Error())
		return mode
	}
	for rows.Next() {
		var m Mode
		rows.Scan(&m.Name, &m.Sales)
		mode = append(mode, m)
	}
	rows.Close()
	return mode
}

func GetMostProfitable(past_days int, n int) []Profit {
	var profit []Profit
	rows, err := conn.Query(context.Background(), "select * from goodbuy.get_top_N_products_by_profit($1, $2)", past_days, n)
	if err != nil {
		println("error getting profit.", err.Error())
		return profit
	}
	for rows.Next() {
		var p Profit
		rows.Scan(&p.Name, &p.Profit)
		profit = append(profit, p)
	}
	rows.Close()
	return profit
}

func GetModeOnMarkets(n int) map[string][]Mode {
	mode := make(map[string][]Mode)
	rows, err := conn.Query(context.Background(), "select * from goodbuy.get_N_popular_products_on_markets($1)", n)
	if err != nil {
		println("error getting mode on markets.", err.Error())
		return mode
	}
	for rows.Next() {
		var m Mode
		var market string
		rows.Scan(&market, &m.Name, &m.Sales)
		mode[market] = append(mode[market], m)
	}
	rows.Close()
	return mode
}
