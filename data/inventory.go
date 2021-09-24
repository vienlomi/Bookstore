package data

import (
	"database/sql"
	"fmt"
)

type Inventory struct {
	InventoryId 		int64  	`json:"inventory_id"`
	Name              	string 	`json:"name"`
	PurchasePrice    	float64	`json:"purchase_price"`
	PurchaseQuantity 	int64  	`json:"purchase_quantity"`
}

type InventoryModel struct {
	DB *sql.DB
}

func (i InventoryModel) Insert(inv *Inventory) error {
	query := `insert into inventory 
			(name, purchase_price, purchase_quantity)
			values (?, ?, ?)`
	_, err := i.DB.Exec(query,
		inv.Name,
		inv.PurchasePrice,
		inv.PurchaseQuantity)
	if err != nil {
		fmt.Println("exec insert fail: ",err)
		return err
	}
	return nil
}

func (i InventoryModel) Update(inv *Inventory) error {
	query := `update inventory
	set name = ?, purchase_price = ?, purchase_quantity = ?
	where inventory_id = ? `
	_, err := i.DB.Exec(query,
		inv.Name,
		inv.PurchasePrice,
		inv.PurchaseQuantity,
		inv.InventoryId)
	if err != nil {
		fmt.Println("exec update fail: ",err)
		return err
	}
	return nil
}

func (i InventoryModel) Delete(id string) error {
	query := `delete from inventory where inventory_id = ?`
	_, err := i.DB.Exec(query, id)
	if err != nil {
		fmt.Println("exec delete fail: ",err)
		return err
	}
	return nil
}

func (i InventoryModel) GetAll() ([]Inventory, error) {
	var inventories []Inventory
	query := `select inventory_id, name, purchase_price, purchase_quantity from inventory`
	rows, err := i.DB.Query(query)
	if err != nil {
		fmt.Println("exec get all fail: ",err)
		return inventories, err
	}
	defer rows.Close()

	for rows.Next() {
		var inv Inventory
		err := rows.Scan(
			&inv.InventoryId,
			&inv.Name,
			&inv.PurchasePrice,
			&inv.PurchaseQuantity)
		if err != nil {
			fmt.Println("scan 1 inventory fail: ",err)
			continue
		}
		inventories = append(inventories, inv)
	}
	return inventories, nil
}

func (i InventoryModel) GetById(id string) (Inventory, error) {
	var inv Inventory
	query := `select inventory_id, name, purchase_price, purchase_quantity from inventory where inventory_id = ?`
	err := i.DB.QueryRow(query, id).Scan(
		&inv.InventoryId,
		&inv.Name,
		&inv.PurchasePrice,
		&inv.PurchaseQuantity)
	if err != nil {
		fmt.Println("get by id inventory fail: ",err)
		return inv, err
	}
	return inv, nil
}
