package db

import (
	"context"
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
	return Product{0, 0, 0, 0, 0, "", ""}
}

func GetProducts(min_threshold, max_threshold Product) []Product {
	rows, err := conn.Query(
		context.Background(),
		"select * from goodbuy.get_filtered_products($1, $2, $3, $4, $5, $6, $7, $8) order by id",
		min_threshold.Category, max_threshold.Category,
		min_threshold.Amount, max_threshold.Amount,
		min_threshold.Self_cost, max_threshold.Self_cost,
		min_threshold.Default_cost, max_threshold.Default_cost)

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
			&product.Category_name,
		)
		if err != nil {
			println("Something on 168", err.Error())
			return nil
		}

		products = append(products, product)
	}
	rows.Close()

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

func GetCategories() []Category {
	rows, err := conn.Query(context.Background(), "select * from goodbuy.product_categories order by id")
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

func AddProductCategory(name string, description string) {
	_, err := conn.Exec(context.Background(), "insert into goodbuy.product_categories values ($1, $2);", name, description)
	if err != nil {
		println("Failed to add a product category", err.Error())
	}
}

func UpdateProductCategory(id uint8, name string, description string) {
	_, err := conn.Exec(context.Background(), "update goodbuy.product_categories set category_name=$1, description=$2 where id=$3", name, description, id)
	if err != nil {
		println("Failed to add a product category", err.Error())
	}
}

func AddProduct(product Product) {
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

func EditProduct(id uint16, new_product Product) {

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
