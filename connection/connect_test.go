package connection

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	Start()

	res, _ := Model.Product.DB.Query("SHOW TABLES")
	defer res.Close()
	var table string
	for res.Next() {
		res.Scan(&table)
		fmt.Println(table)
	}
}
