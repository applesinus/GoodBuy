package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn
var connection_err error

type Product struct {
	Category      int8
	Amount        uint16
	Id            uint16
	Self_cost     float32
	Default_cost  float32
	Name          string
	Category_name string
}

type Category struct {
	Cat_id          int8
	Cat_description string
	Cat_name        string
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
		println("Something on 83", err.Error())
		return false
	} else {
		password := ""
		err := conn.QueryRow(context.Background(), "select password from users where id=$1", id).Scan(&password)
		if err != nil {
			println("Something on 89", err.Error())
			return false
		} else {
			if password == "" || password != inputed_password {
				println("wrong password:", inputed_password, "expected:", password)
				return false
			}
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

	select_contitions.Grow(len(" order by id"))
	select_contitions.WriteString(" order by id")

	conditions := select_contitions.String()
	rows, err := conn.Query(context.Background(), conditions)

	if err != nil {
		println("Something on 150", err.Error())
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
			println("Something on 168", err.Error())
			return nil
		}

		products = append(products, product)
	}
	rows.Close()

	for i := 0; i < len(products); i++ {
		var cat_name string
		err := conn.QueryRow(context.Background(), "select category_name from product_categories where id=$1", products[i].Category).Scan(&cat_name)
		if err != nil {
			println("Something on 180", err.Error())
		}
		products[i].Category_name = cat_name
	}

	return products
}

func GetProductByID(id int) Product {
	prod, _ := conn.Query(context.Background(), "select * from products where id=$1", id)
	product := NewProduct()
	for prod.Next() {
		prod.Scan(
			&product.Name,
			&product.Default_cost,
			&product.Category,
			&product.Self_cost,
			&product.Amount,
			&product.Id,
		)
	}
	return product
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
		println("Something on 239", err.Error())
		return nil
	}

	categories := make([]Category, 0)

	for rows.Next() {
		var cat Category

		err = rows.Scan(
			&cat.Cat_name,
			&cat.Cat_description,
			&cat.Cat_id,
		)
		if err != nil {
			println("Something on 254", err.Error())
			return nil
		}

		categories = append(categories, cat)
	}
	rows.Close()

	return categories
}

func GetCategoryIdByName(name string) int8 {
	var id int8
	err := conn.QueryRow(context.Background(), "select id from product_categories where category_name=$1", name).Scan(&id)
	if err != nil {
		println("Something on 269", err.Error())
		return 0
	}
	return id
}

func GetCategoryNameById(id int8) string {
	var name string
	err := conn.QueryRow(context.Background(), "select category_name from product_categories where id=$1", id).Scan(&name)
	if err != nil {
		println("Something on 269", err.Error())
		err = conn.QueryRow(context.Background(), "select category_name from product_categories where id=0").Scan(&name)
		if err != nil {
			return "ОШИБКА В БАЗЕ ДАННЫХ"
		}
	}
	return name
}

func AddProduct(product Product) {
	var query strings.Builder
	var val string
	query.Grow(len("insert into products values ('"))
	query.WriteString("insert into products values ('")

	query.Grow(len(product.Name))
	query.WriteString(product.Name)
	query.Grow(len("', "))
	query.WriteString("', ")

	val = fmt.Sprintf("%.2f", product.Default_cost)
	query.Grow(len(val))
	query.WriteString(val)
	query.Grow(len(", "))
	query.WriteString(", ")

	val = fmt.Sprintf("%v", product.Category)
	query.Grow(len(val))
	query.WriteString(val)
	query.Grow(len(", "))
	query.WriteString(", ")

	val = fmt.Sprintf("%.2f", product.Self_cost)
	query.Grow(len(val))
	query.WriteString(val)
	query.Grow(len(", "))
	query.WriteString(", ")

	val = fmt.Sprintf("%v", product.Amount)
	query.Grow(len(val))
	query.WriteString(val)
	query.Grow(len(")"))
	query.WriteString(")")

	println(query.String())

	conn.Exec(context.Background(), query.String())
}

func EditProduct(id uint16, new_product Product) {
	var val string
	var query strings.Builder
	query.Grow(len("update products set name = '"))
	query.WriteString("update products set name = '")

	query.Grow(len(new_product.Name))
	query.WriteString(new_product.Name)
	query.Grow(len("', default_cost = "))
	query.WriteString("', default_cost = ")

	val = fmt.Sprintf("%.2f", new_product.Default_cost)
	query.Grow(len(val))
	query.WriteString(val)
	query.Grow(len(", category = "))
	query.WriteString(", category = ")

	val = fmt.Sprintf("%v", new_product.Category)
	query.Grow(len(val))
	query.WriteString(val)
	query.Grow(len(", self_cost = "))
	query.WriteString(", self_cost = ")

	val = fmt.Sprintf("%.2f", new_product.Self_cost)
	query.Grow(len(val))
	query.WriteString(val)
	query.Grow(len(", amount = "))
	query.WriteString(", amount = ")

	val = fmt.Sprintf("%v", new_product.Amount)
	query.Grow(len(val))
	query.WriteString(val)

	query.Grow(len(" where id = "))
	query.WriteString(" where id = ")
	val = fmt.Sprintf("%v", id)
	query.Grow(len(val))
	query.WriteString(val)

	conn.Exec(context.Background(), query.String())
}
