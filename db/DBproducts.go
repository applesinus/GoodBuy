package db

import (
	"context"
	"strconv"
	"strings"
)

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

func NewProduct() Product {
	return Product{0, 0, 0, 0, 0, "0", ""}
}

// TODO refactor NOW
func GetProducts(min_threshold, max_threshold Product) []Product {
	var select_contitions strings.Builder
	select_contitions.Grow(len("select * from goodbuy.products"))
	select_contitions.WriteString("select * from goodbuy.products")

	// Category trashold
	new_cond := new_condition(
		select_contitions.Len() != len("select * from goodbuy.products"),
		"category",
		strconv.Itoa(int(min_threshold.Category)),
		strconv.Itoa(int(max_threshold.Category)),
	)
	select_contitions.Grow(len(new_cond))
	select_contitions.WriteString(new_cond)

	// Amount trashold
	new_cond = new_condition(
		select_contitions.Len() != len("select * from goodbuy.products"),
		"amount",
		strconv.Itoa(int(min_threshold.Amount)),
		strconv.Itoa(int(max_threshold.Amount)),
	)
	select_contitions.Grow(len(new_cond))
	select_contitions.WriteString(new_cond)

	// Self_cost trashold
	new_cond = new_condition(
		select_contitions.Len() != len("select * from goodbuy.products"),
		"self_cost",
		strconv.Itoa(int(min_threshold.Self_cost)),
		strconv.Itoa(int(max_threshold.Self_cost)),
	)
	select_contitions.Grow(len(new_cond))
	select_contitions.WriteString(new_cond)

	// Default_cost trashold
	new_cond = new_condition(
		select_contitions.Len() != len("select * from goodbuy.products"),
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
		err := conn.QueryRow(context.Background(), "select category_name from goodbuy.product_categories where id=$1", products[i].Category).Scan(&cat_name)
		if err != nil {
			println("Something on 180", err.Error())
		}
		products[i].Category_name = cat_name
	}

	return products
}

func GetProductByID(id int) Product {
	prod, _ := conn.Query(context.Background(), "select * from goodbuy.products where id=$1", id)
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

func GetCategories() []Category {
	rows, err := conn.Query(context.Background(), "select * from goodbuy.product_categories")
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
	err := conn.QueryRow(context.Background(), "select id from goodbuy.product_categories where category_name=$1", name).Scan(&id)
	if err != nil {
		println("Something on 269", err.Error())
		return 0
	}
	return id
}

func GetCategoryNameById(id int8) string {
	var name string
	err := conn.QueryRow(context.Background(), "select category_name from goodbuy.product_categories where id=$1", id).Scan(&name)
	if err != nil {
		println("Something on 269", err.Error())
		err = conn.QueryRow(context.Background(), "select category_name from goodbuy.product_categories where id=0").Scan(&name)
		if err != nil {
			return "ОШИБКА В БАЗЕ ДАННЫХ"
		}
	}
	return name
}

// refactor done
func AddProduct(product Product) {
	/*var query strings.Builder
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

	conn.Exec(context.Background(), query.String())*/

	_, err := conn.Exec(
		context.Background(),
		"call goodbuy.add_product($1, $2, $3, $4, $5)",
		product.Name,
		product.Default_cost,
		product.Category_name,
		product.Self_cost,
		product.Amount,
	)
	if err != nil {
		println("can't add product", err.Error())
	}
}

// refactor done
func EditProduct(id uint16, new_product Product) {
	/*var val string
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

	conn.Exec(context.Background(), query.String())*/

	_, err := conn.Exec(context.Background(),
		"call goodbuy.edit_product($1, $2, $3, $4, $5, $6)",
		id,
		new_product.Name,
		new_product.Default_cost,
		new_product.Category_name,
		new_product.Self_cost,
		new_product.Amount,
	)
	if err != nil {
		println("can't edit the product, id:", id, err.Error())
	}
}
