package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn
var connection_err error

var psql_port = "10000/postgres"

func execute_file(file_name string) {
	commands, err := os.ReadFile(file_name)
	if err != nil {
		println("NO SUCH FILE!")
		return
	}

	_, err = conn.Exec(context.Background(), string(commands))
	if err != nil {
		println("An error in file ", file_name, ": ", err.Error())
	}
}

func StartDB() {
	var err error

	conn, connection_err = pgx.Connect(context.Background(), "postgres://user:passw0rd@localhost:"+psql_port)
	if connection_err != nil {
		println("DB CONNECTION IS FAILED")
		panic(connection_err)
	}

	// test if the DB exist or corrupted. if this is the case for now it fully remakes (!not recover!) the DB
	var test string
	err = conn.QueryRow(context.Background(), "select description from goodbuy.statuses where id=$1", 1).Scan(&test)
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
	execute_file("db/CreateDB.sql")
	println("DB created")
}

func DropDB() {
	execute_file("db/DropDB.sql")
	println("DB dropped")
}

func AddTestData() {
	execute_file("db/TestData.sql")
	println("Test data added")
}

func Auth(inputed_username, inputed_password string) bool {
	id := -1
	err := conn.QueryRow(context.Background(), "select id from goodbuy.users where username=$1", inputed_username).Scan(&id)
	if err != nil {
		return false
	} else {
		password := ""
		err := conn.QueryRow(context.Background(), "select password from goodbuy.users where id=$1", id).Scan(&password)
		if err != nil {
			return false
		} else {
			if password == "" || inputed_password != password {
				return false
			}
			return true
		}
	}
}
