package data

import "database/sql"

type Models struct {
	Product    	ProductModel
	Inventory  	InventoryModel
	User       	UserModel
	Order		OrderModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Product:    ProductModel{DB: db},
		Inventory:  InventoryModel{DB: db},
		User:       UserModel{DB: db},
		Order: 		OrderModel{DB: db},
	}
}
