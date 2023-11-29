package main

import (
	"context"
	"os"
)

func execute_file_SQL_sequence(file_name string) {
	commands, err := os.ReadFile(file_name)
	if err != nil {
		println("NO SUCH FILE!")
		return
	}
	conn.Exec(context.Background(), string(commands))
}

func CreateDB() {
	execute_file_SQL_sequence("PostgreSQL/CreateDB.txt")
}

func DropDB() {
	execute_file_SQL_sequence("PostgreSQL/DropDB.txt")
}
