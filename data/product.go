package data

import (
	"database/sql"
	"errors"
	"fmt"
)

type Product struct {
	ProductId     int64   `json:"product_id"`
	Name          string  `json:"name"`
	Author        string  `json:"author"`
	Collection    string  `json:"collection"`
	ImageUrl      string  `json:"image_url"`
	Price         float64 `json:"price"`
	Sales         float64 `json:"sales"`
	Publisher     string  `json:"publisher"`
	DatePublish   string  `json:"date_publish"`
	Rate          float64 `json:"rate"`
	Description   string  `json:"description"`
	StatusProduct float64 `json:"status"`
}

type ProductLite struct {
	ProductId  int64   `json:"product_id"`
	Name       string  `json:"name"`
	Author     string  `json:"author"`
	Collection string  `json:"collection"`
	ImageUrl   string  `json:"image_url"`
	Price      float64 `json:"price"`
	Sales      float64 `json:"sales"`
}

type ProductModel struct {
	DB *sql.DB
}

func (p *ProductModel) Insert(product *Product) error {
	fmt.Println(product)
	query := `insert into products 
	(name, author, collection, image_url, price, sales, publisher, date_publish, rate, description, status_product)
	values (?, ?, ?, ?, ?, ? ,?, ?, ?, ?, ?)`

	_, err := p.DB.Exec(query,
		product.Name,
		product.Author,
		product.Collection,
		product.ImageUrl,
		product.Price,
		product.Sales,
		product.Publisher,
		product.DatePublish,
		product.Rate,
		product.Description,
		product.StatusProduct)
	if err != nil {
		fmt.Println("exec insert fail: ", err)
		return err
	}
	return nil
}

func (p *ProductModel) Update(product *Product) error {
	query := `update products
	set name = ?, author = ?, collection = ?, image_url = ?, price = ?, sales = ?, publisher = ?, 
	    date_publish = ?, rate = ?, description = ?, status_product = ? where product_id = ?`
	row, err := p.DB.Exec(query,
		product.Name,
		product.Author,
		product.Collection,
		product.ImageUrl,
		product.Price,
		product.Sales,
		product.Publisher,
		product.DatePublish,
		product.Rate,
		product.Description,
		product.StatusProduct,
		product.ProductId)
	if err != nil {
		fmt.Println("exec update fail: ", err)
		return err
	}
	if stt, _ := row.RowsAffected(); stt == 0 {
		return errors.New("id not available")
	}
	return nil
}

func (p *ProductModel) Delete(id string) error {
	query := `delete from products where product_id = ?`
	row, err := p.DB.Exec(query, id)
	if err != nil {
		fmt.Println("exec delete fail: ", err)
		return err
	}
	stt, _ := row.RowsAffected()
	if stt == 0 {
		return errors.New("id not available")
	}
	return nil
}

func (p *ProductModel) GetDetailById(id string) (Product, error) {
	product := Product{}
	query := `select product_id, name, author, collection, image_url, price, sales, publisher,
       date_publish, rate, description, status_product from products where product_id = ? `
	err := p.DB.QueryRow(query, id).Scan(
		&product.ProductId,
		&product.Name,
		&product.Author,
		&product.Collection,
		&product.ImageUrl,
		&product.Price,
		&product.Sales,
		&product.Publisher,
		&product.DatePublish,
		&product.Rate,
		&product.Description,
		&product.StatusProduct)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (p *ProductModel) GetAll(payload Payload) ([]ProductLite, error) {
	offset := (payload.Stt - 1) * payload.Limit // bat dau la row 0
	var products []ProductLite
	query := `select product_id, name, author, collection, image_url, price, sales from products LIMIT ?, ?`
	rows, err := p.DB.Query(query, offset, payload.Limit)
	if err != nil {
		fmt.Println("exec get all product fail: ", err)
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product ProductLite
		err := rows.Scan(
			&product.ProductId,
			&product.Name,
			&product.Author,
			&product.Collection,
			&product.ImageUrl,
			&product.Price,
			&product.Sales)
		if err != nil {
			fmt.Println("scan 1 product lite fail ", err)
			continue
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductModel) GetListProduct(payload Payload, orderBy string) ([]ProductLite, error) {
	offset := (payload.Stt - 1) * payload.Limit // bat dau la row 0
	var products []ProductLite
	queryText := fmt.Sprintf("select product_id, name, author, collection, image_url, price, sales from products ORDER BY %s DESC LIMIT %d, %d", orderBy, offset, payload.Limit)
	rows, err := p.DB.Query(queryText)
	if err != nil {
		fmt.Println("exec get latest products fail: ", err)
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product ProductLite
		err := rows.Scan(
			&product.ProductId,
			&product.Name,
			&product.Author,
			&product.Collection,
			&product.ImageUrl,
			&product.Price,
			&product.Sales)
		if err != nil {
			fmt.Println("scan 1 product lite fail ", err)
			continue
		}
		products = append(products, product)
	}
	return products, nil
}
