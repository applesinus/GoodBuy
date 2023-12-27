package db

import "context"

func GetAllIncome() int {
	sum := 0
	// TODO
	// это - заглушка для зачета по го, чтобы показать функциональность.
	// не смотрите на реализацию с точки зрения БД
	// запрос в красивом виде в tmpTotal.sql
	err := conn.QueryRow(context.Background(), "select sum(p.count * p.cost) as total from positions p join positions_in_receipts pr on p.id = pr.position join receipts r on pr.receipt = r.id where r.status = 1 and p.status = 1;").Scan(&sum)
	if err != nil {
		println("error getting total.", err.Error())
	}
	return sum
}
