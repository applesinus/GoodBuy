package db

import "github.com/jackc/pgtype"

type Position struct {
	Status     uint8
	Count      uint8
	Id         uint16
	Cost       float32
	Product_id uint16
}

type Receipt struct {
	Status       uint8
	Id           uint16
	Total        uint16
	Positions_id []uint16
	Date         pgtype.Date
}
