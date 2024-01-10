package db

import (
	"context"
	"time"
)

type Position struct {
	Count   uint8
	Id      uint16
	Cost    float32
	Product string
	Status  string
}

type Receipt struct {
	Id        uint16
	Total     float64
	Pos_len   int
	Positions []Position
	Status    string
	Date      string
}

func NewPosition() Position {
	return Position{0, 0, 0, "OK", ""}
}

func NewReceipt() Receipt {
	pos := make([]Position, 0)
	pos_len := len(pos) + 1
	return Receipt{0, 0, pos_len, pos, "OK", "1999-01-01"}
}

func GetAllReceipts() map[int]Receipt {

	receipts := make(map[int]Receipt)

	rows, err := conn.Query(context.Background(), "select * from goodbuy.get_detailed_receipts")

	if err != nil {
		println("err on getting receipts", err.Error())
	}

	for rows.Next() {
		position := NewPosition()
		receipt := NewReceipt()
		var date time.Time
		err := rows.Scan(
			&date,
			&receipt.Id,
			&receipt.Status,
			&position.Id,
			&position.Product,
			&position.Cost,
			&position.Count,
			&position.Status,
		)
		receipt.Date = date.Format("2006-01-02")
		if err != nil {
			println("err on setting reciept's values", err.Error())
			return receipts
		}
		id := int(receipt.Id)

		if _, ok := receipts[int(receipt.Id)]; !ok {
			receipts[id] = receipt
		}

		total := float64(position.Cost) * float64(position.Count)
		currentReceipt := receipts[id]
		currentReceipt.Total += total
		currentReceipt.Pos_len += 1
		currentReceipt.Positions = append(currentReceipt.Positions, position)
		receipts[id] = currentReceipt
	}
	rows.Close()

	return receipts
}
