package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func execute_file(file_name string) {
	commands, err := os.ReadFile(file_name)
	if err != nil {
		println("NO SUCH FILE!")
		return
	}
	conn.Exec(context.Background(), string(commands))
}

func StartDB() {
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

func CreateDB() {
	execute_file("db/CreateDB.txt")
	println("DB created")
}

func DropDB() {
	execute_file("db/DropDB.txt")
	println("DB created")
}
