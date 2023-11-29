package main

import (
	"context"

	"github.com/jackc/pgx"
)

var conn *pgx.Conn

func main() {
	var err error
	conn, _ = pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/postgres")

	// test if the DB exist or corrupted. if this is the case for now it fully remakes (!not recover!) the DB
	var test string
	err = conn.QueryRow(context.Background(), "select description from statuses where id=$1", 1).Scan(&test)
	if err != nil {
		// DB doesn't seem to exist
		CreateDB()
	} else if test != "OK" {
		// DB doesn't seem to be valid
		DropDB()
		CreateDB()
	}
}
