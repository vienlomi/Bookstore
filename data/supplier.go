package data

import (
	"database/sql"
	"fmt"
)

type Supplier struct {
	SupplierId  int64		`json:"supplier_id"`
	CompanyName string		`json:"company_name"`
	Address     string		`json:"address"`
}

type SupplierModel struct {
	DB *sql.DB
}

func (s SupplierModel) Insert(sup *Supplier) error {
	query := `insert into suppliers (company_name, address) values (?, ?)`
	_, err := s.DB.Exec(query, sup.CompanyName, sup.Address)
	if err != nil {
		fmt.Println("exec insert sup fail ",err)
		return err
	}
	return nil
}

func (s SupplierModel) Update(sup *Supplier) error {
	query := `update suppliers 
		set company_name = ?, address= ? where supplier_id = ?`
	_, err := s.DB.Exec(query, sup.CompanyName, sup.Address, sup.SupplierId)
	return err
}
func (s SupplierModel) Delete(id string, sup *Supplier) error {
	query := `delete from suppliers where supplier_id = ?`
	_, err := s.DB.Exec(query, id)
	if err != nil {
		fmt.Println("exec update sup fail ",err)
		return err
	}
	return nil
}
func (s SupplierModel) GetById(id string) (Supplier, error) {
	var sup Supplier
	query := `select supplier_id, company_name, address from suppliers where supplier_id = ?`
	err := s.DB.QueryRow(query, id).Scan(
		&sup.SupplierId,
		&sup.CompanyName,
		&sup.Address)
	if err != nil {
		fmt.Println("exec update sup fail ",err)
		return sup, err
	}
	return sup, nil
}
func (s SupplierModel) GetAll() ([]Supplier, error) {
	var suppliers []Supplier
	query := `select supplier_id, company_name, address from supplier`
	rows, err := s.DB.Query(query)
	if err != nil {
		return suppliers, err
	}
	for rows.Next() {
		var sup Supplier
		err := rows.Scan(
			&sup.SupplierId,
			&sup.CompanyName,
			&sup.Address)
		if err != nil {
			fmt.Println("scan 1 supplier fail")
			continue
		}
		suppliers = append(suppliers, sup)
	}
	return suppliers, nil
}
