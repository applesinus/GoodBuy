package db

import (
	"context"
	"time"
)

type Position struct {
	Status  uint8
	Count   uint8
	Id      uint16
	Cost    float32
	Product string
}

type Receipt struct {
	Status    uint8
	Id        uint16
	Total     uint16
	Pos_len   int
	Positions []Position
	Date      string
}

func NewPosition() Position {
	return Position{0, 0, 0, 0, ""}
}

func NewReceipt() Receipt {
	pos := make([]Position, 0)
	pos_len := len(pos) + 1
	return Receipt{0, 0, 0, pos_len, pos, "1999-01-01"}
}

func GetAllReceipts() map[int]Receipt {

	// getting all receipts
	receipts := make(map[int]Receipt)
	rows, err := conn.Query(context.Background(), "select * from receipts")
	if err != nil {
		println("err on getting receipts", err.Error())
	}
	for rows.Next() {
		reciept := NewReceipt()
		var date time.Time
		err := rows.Scan(
			&date,
			&reciept.Status,
			&reciept.Id,
		)
		if err != nil {
			println("err on setting reciept's values", err.Error())
			return nil
		}
		reciept.Date = date.Format("2006-01-02")
		receipts[int(reciept.Id)] = reciept
	}
	rows.Close()

	// getting all positions
	positions := make(map[int]Position)
	rows, err = conn.Query(context.Background(), "select * from positions")
	if err != nil {
		println("err on getting positions", err.Error())
	}
	for rows.Next() {
		pos := NewPosition()
		var prod uint16
		err := rows.Scan(
			&prod,
			&pos.Cost,
			&pos.Count,
			&pos.Status,
			&pos.Id,
		)
		if err != nil {
			println("err on setting position's values", err.Error())
			return nil
		}
		pos.Product = GetProductByID(int(prod)).Name
		positions[int(pos.Id)] = pos
	}
	rows.Close()

	// adding positions to receipts
	rows, err = conn.Query(context.Background(), "select * from positions_in_receipts")
	if err != nil {
		println("err on getting pos_in_rec", err.Error())
	}
	for rows.Next() {
		var pos, rec int
		err := rows.Scan(
			&pos,
			&rec,
		)
		if err != nil {
			println("err on setting pos_in_rec's values", err.Error())
			return nil
		}
		reciept := receipts[rec]
		reciept.Pos_len++
		reciept.Positions = append(reciept.Positions, positions[pos])
		receipts[rec] = reciept
	}
	rows.Close()

	return receipts
}
