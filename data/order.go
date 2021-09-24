package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Order struct {
	OrderId         int64   `json:"order_id"`
	UserId          int64   `json:"user_id"`
	TotalPrice      float64 `json:"total_price"`
	ShippingMethod  string  `json:"shipping_method"`
	ReceiverAddress string  `json:"receive_address"`
	ReceiverPhone   string  `json:"receive_phone"`
	Note            string  `json:"note"`
	PayMethod       string  `json:"pay_method"`
	NumberCart      string  `json:"number_cart"`
	OwnerName       string  `json:"owner_name"`
	ShippedDate     string  `json:"shipped_date"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`

	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductId int64   `json:"product_id"`
	Name      string  `json:"name"`
	ImageUrl  string  `json:"image_url"`
	Price     float64 `json:"price"`
	Amount    int     `json:"amount"`
	UnitPrice float64 `json:"unit_price"`
}

type OrderLite struct {
	OrderId    int64   `json:"order_id"`
	UserID     int64   `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	CreateAt   string  `json:"create_at"`
}

type OrderModel struct {
	DB *sql.DB
}

func (o *Order) String() string {
	str, err := json.Marshal(o)
	if err != nil {
		fmt.Println("can not marshal order")
		return ""
	}
	return string(str)
}

func (o *OrderModel) Insert(order *Order) error {
	tx, err := o.DB.Begin()
	if err != nil {
		fmt.Println("insert DB begin fail: ", err)
		return err
	}
	defer tx.Rollback()

	query := `insert into orders (user_id, total_price, shipping_method, receiver_address, receiver_phone, 
	note, pay_method, number_cart, owner_name) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	rs, err := tx.Exec(query, order.UserId, order.TotalPrice, order.ShippingMethod, order.ReceiverAddress, order.ReceiverPhone,
		order.Note, order.PayMethod, order.NumberCart, order.OwnerName)
	if err != nil {
		fmt.Println("exec insert 1 fail: ", err)
		return err
	}
	order.OrderId, err = rs.LastInsertId()
	if err != nil {
		fmt.Println("get last id insert 1 fail: ", err)
		return err
	}

	sttm, err := tx.Prepare(`insert into order_items (order_id, product_id, quantity, unit_price) values (?, ?, ?, ?)`)
	if err != nil {
		fmt.Println("prepare insert 2 fail: ", err)
		return err
	}
	defer sttm.Close()

	for _, item := range order.Items {
		if _, err := sttm.Exec(order.OrderId, item.ProductId, item.Amount, item.UnitPrice); err != nil {
			fmt.Println("exec insert 2 fail", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("commit insert order fail ", err)
		return err
	}
	return nil
}

func (o *OrderModel) Delete(id int) error {
	query := `delete from orders where order_id = ?`
	_, err := o.DB.Query(query, id)
	if err != nil {
		fmt.Println("delete order fail ", err)
		return err
	}
	return nil
}

func (o *OrderModel) Update(id int, shippedDate time.Duration, status string) error {
	query := `update orders set shipped_date = ?, status = ? where order_id =?`
	_, err := o.DB.Query(query, shippedDate, status, id)
	if err != nil {
		fmt.Println("update order fail ", err)
		return err
	}
	return nil
}

func (o *OrderModel) GetAllASC(payload Payload) ([]OrderLite, error) {
	offset := (payload.Stt - 1) * payload.Limit // bat dau la row 0
	var listOrder []OrderLite
	query := `select order_id, user_id, total_price, IFNULL(status, ""), created_at from orders ORDER BY ASC LIMIT ?, ?`
	rows, err := o.DB.Query(query, offset, payload.Limit)
	if err != nil {
		fmt.Println("get all orders fail ", err)
		return listOrder, err
	}
	for rows.Next() {
		var order OrderLite
		err = rows.Scan(&order.OrderId, &order.UserID, &order.TotalPrice, &order.Status, &order.CreateAt)
		if err != nil {
			fmt.Println("scan 1 order lite fail ", err)
			continue
		}
		listOrder = append(listOrder, order)
	}
	return listOrder, nil
}

func (o *OrderModel) GetAllDESC(payload Payload) ([]OrderLite, error) {
	offset := (payload.Stt - 1) * payload.Limit // bat dau la row 0
	var listOrder []OrderLite
	queryText := fmt.Sprintf("select order_id, user_id, total_price, IFNULL(status, ''), created_at from orders ORDER BY created_at DESC LIMIT %d, %d", offset, payload.Limit)
	rows, err := o.DB.Query(queryText)
	if err != nil {
		fmt.Println("get all orders fail ", err)
		return listOrder, err
	}
	for rows.Next() {
		var order OrderLite
		err = rows.Scan(&order.OrderId, &order.UserID, &order.TotalPrice, &order.Status, &order.CreateAt)
		if err != nil {
			fmt.Println("scan 1 order lite fail ", err)
			continue
		}
		listOrder = append(listOrder, order)
	}
	return listOrder, nil
}

func (o *OrderModel) GetDetailOrderBuyId(id string) (Order, error) {
	var order Order
	var item OrderItem
	query := `SELECT c.order_id, c.user_id, c.total_price, IFNULL(c.status, ""), c.created_at, c.quantity, c.unit_price,
        d.product_id, d.name, d.image_url, d.price FROM
(select a.order_id, a.user_id, a.total_price, a.status, a.created_at, b.product_id, b.quantity, b.unit_price from
(select * from orders where user_id = ?) as a
INNER JOIN
(select * from order_items) as b
WHERE a.order_id = b.order_id
) as c
INNER JOIN
(select product_id, name, image_url, price from products) as d
WHERE c.product_id = d.product_id;`

	rows, err := o.DB.Query(query, id)
	if err != nil {
		fmt.Println("get order buy id fail ", err)
		return order, err
	}
	//var item OrderItem
	for rows.Next() {

		err = rows.Scan(&order.OrderId, &order.UserId, &order.TotalPrice, &order.Status, &order.CreatedAt,
			&item.Amount, &item.UnitPrice, &item.ProductId, &item.Name, &item.ImageUrl, &item.Price)
		if err != nil {
			fmt.Println("scan db fail ", err)
			return order, err
		}
		itemCopy := item
		order.Items = append(order.Items, itemCopy)
		//order.Items = append(order.Items, item)
	}
	return order, nil
}

func (o *OrderModel) GetOrdersByUserId(id string, payload Payload) ([]Order, error) {

	var listOrder []Order
	offset := (payload.Stt - 1) * payload.Limit // bat dau la row 0
	queryText := fmt.Sprintf("SELECT c.order_id, c.total_price, IFNULL(c.status, \"\"), c.created_at, c.quantity, c.unit_price,\n        d.product_id, d.name, d.image_url, d.price FROM\n(select a.order_id, a.user_id, a.total_price, a.status, a.created_at, b.product_id, b.quantity, b.unit_price from\n(select * from orders where user_id = %s LIMIT %d, %d) as a\nINNER JOIN\n(select * from order_items) as b\nWHERE a.order_id = b.order_id\n) as c\nINNER JOIN\n(select product_id, name, image_url, price from products) as d\nWHERE c.product_id = d.product_id\nORDER BY c.created_at DESC", id, offset, payload.Limit)
	rows, err := o.DB.Query(queryText)
	if err != nil {
		fmt.Println("get orders buy user_id fail ", err)
		return nil, err
	}
	for rows.Next() {
		var order Order
		var item OrderItem
		err = rows.Scan(&order.OrderId, &order.TotalPrice, &order.Status, &order.CreatedAt, &item.Amount, &item.UnitPrice, &item.ProductId, &item.Name, &item.ImageUrl, &item.Price)
		if err != nil {
			fmt.Println("scan 1 order detail fail ", err)
			continue
		}
		if len(listOrder) == 0 {
			order.Items = append(order.Items, item)
			listOrder = append(listOrder, order)
		} else if listOrder[len(listOrder)-1].OrderId != order.OrderId {
			order.Items = append(order.Items, item)
			listOrder = append(listOrder, order)
		} else {
			listOrder[len(listOrder)-1].Items = append(listOrder[len(listOrder)-1].Items, item)
		}

	}
	return listOrder, nil
}

type DetailOrder struct {
	Id       int64   `json:"id"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
	Image    string  `json:"image"`
	Name     string  `json:"name"`
}

func (o *OrderModel) GetOrders(id string) ([]DetailOrder, error) {
	var listOrder []DetailOrder
	query := `select ord.order_id, item.quantity, p.price, p.image_url, p.name from orders as ord join order_items as item using (order_id) join products as p using (product_id) where ord.user_id = ?`
	rows, err := o.DB.Query(query, id)
	if err != nil {
		log.Println(err)
		return listOrder, err
	}
	for rows.Next() {
		var order DetailOrder
		err = rows.Scan(
			&order.Id,
			&order.Quantity,
			&order.Price,
			&order.Image,
			&order.Name,
		)
		listOrder = append(listOrder, order)
	}
	return listOrder, nil
}
