package redis

import (
	"data"
	"fmt"
	"testing"
)

func TestGetList(t *testing.T) {
	NewClient()
	var newProducts []data.ProductLite
	err := ClientLRangProduct(data.OutStanding["1"], &newProducts)
	if err != nil {
		t.Error("can not get list")
	}
	fmt.Println(newProducts)
}
