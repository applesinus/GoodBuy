package db

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn
var connection_err error

type Product struct {
	Category      int8
	Amount        int8
	Self_cost     float32
	Default_cost  float32
	Id            int
	Name          string
	Category_name string
}

type Category struct {
	Id          int8
	Description string
	Name        string
}

var psql_port = "10000/postgres"

func execute_file(file_name string) {
	commands, err := os.ReadFile(file_name)
	if err != nil {
		println("NO SUCH FILE!")
		return
	}

	_, err = conn.Exec(context.Background(), string(commands))
	if err != nil {
		println("An error in tasked file ", file_name, ": ", err.Error())
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
	execute_file("db/CreateDB.sql")
	println("DB created")
}

func DropDB() {
	execute_file("db/DropDB.sql")
	println("DB dropped")
}

func Auth(inputed_username, inputed_password string) bool {
	id := -1
	err := conn.QueryRow(context.Background(), "select id from users where username=$1", inputed_username).Scan(&id)
	if err != nil {
		println()
		println(err.Error())
		return false
	} else {
		password := ""
		err := conn.QueryRow(context.Background(), "select password from users where id=$1", id).Scan(&password)
		if err != nil || password == "" || password != inputed_password {
			println()
			println(err.Error())
			return false
		} else {
			return true
		}
	}
}

func GetProducts(min_threshold, max_threshold Product) []Product {
	var select_contitions strings.Builder
	select_contitions.Grow(len("select * from products"))
	select_contitions.WriteString("select * from products")

	// Category trashold
	new_cond := new_condition(
		select_contitions.Len() != len("select * from products"),
		"category",
		strconv.Itoa(int(min_threshold.Category)),
		strconv.Itoa(int(max_threshold.Category)),
	)
	select_contitions.Grow(len(new_cond))
	select_contitions.WriteString(new_cond)

	// Amount trashold
	new_cond = new_condition(
		select_contitions.Len() != len("select * from products"),
		"amount",
		strconv.Itoa(int(min_threshold.Amount)),
		strconv.Itoa(int(max_threshold.Amount)),
	)
	select_contitions.Grow(len(new_cond))
	select_contitions.WriteString(new_cond)

	// Self_cost trashold
	new_cond = new_condition(
		select_contitions.Len() != len("select * from products"),
		"self_cost",
		strconv.Itoa(int(min_threshold.Self_cost)),
		strconv.Itoa(int(max_threshold.Self_cost)),
	)
	select_contitions.Grow(len(new_cond))
	select_contitions.WriteString(new_cond)

	// Default_cost trashold
	new_cond = new_condition(
		select_contitions.Len() != len("select * from products"),
		"default_cost",
		strconv.Itoa(int(min_threshold.Default_cost)),
		strconv.Itoa(int(max_threshold.Default_cost)),
	)
	select_contitions.Grow(len(new_cond))
	select_contitions.WriteString(new_cond)

	conditions := select_contitions.String()
	rows, err := conn.Query(context.Background(), conditions)

	if err != nil {
		println(err.Error())
		return nil
	}

	products := make([]Product, 0)

	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.Name,
			&product.Default_cost,
			&product.Category,
			&product.Self_cost,
			&product.Amount,
			&product.Id,
		)
		if err != nil {
			println(err.Error())
			return nil
		}

		products = append(products, product)
	}
	rows.Close()

	for i := 0; i < len(products); i++ {
		var cat_name string
		err := conn.QueryRow(context.Background(), "select category_name from product_categories where id=$1", products[i].Category).Scan(&cat_name)
		if err != nil {
			println(err.Error())
		}
		products[i].Category_name = cat_name
	}

	return products
}

func new_condition(notFirst bool, value, min, max string) string {
	var new_cond strings.Builder

	if notFirst {
		new_cond.Grow(len(" and "))
		new_cond.WriteString(" and ")
	} else {
		new_cond.Grow(len(" where "))
		new_cond.WriteString(" where ")
	}

	if max != "0" {
		new_cond.Grow(len(value))
		new_cond.WriteString(value)
		new_cond.Grow(len(" <= "))
		new_cond.WriteString(" <= ")
		new_cond.Grow(len(max))
		new_cond.WriteString(max)
		if min != "0" {
			new_cond.Grow(len(" and "))
			new_cond.WriteString(" and ")
			new_cond.Grow(len(value))
			new_cond.WriteString(value)
			new_cond.Grow(len(" >= "))
			new_cond.WriteString(" >= ")
			new_cond.Grow(len(min))
			new_cond.WriteString(min)
		}
	} else {
		if min != "0" {
			new_cond.Grow(len(value))
			new_cond.WriteString(value)
			new_cond.Grow(len(" >= "))
			new_cond.WriteString(" >= ")
			new_cond.Grow(len(min))
			new_cond.WriteString(min)
		} else {
			return ""
		}
	}

	return new_cond.String()

}

func NewProduct() Product {
	return Product{0, 0, 0, 0, 0, "0", ""}
}

func GetCategories() []Category {
	rows, err := conn.Query(context.Background(), "select * from product_categories")
	if err != nil {
		println(err.Error())
		return nil
	}

	categories := make([]Category, 0)

	for rows.Next() {
		var cat Category

		err = rows.Scan(
			&cat.Name,
			&cat.Description,
			&cat.Id,
		)
		if err != nil {
			println(err.Error())
			return nil
		}

		categories = append(categories, cat)
	}
	rows.Close()

	return categories
}
