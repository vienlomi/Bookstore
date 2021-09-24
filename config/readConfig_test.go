package config

import (
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	ReadConfigCsv("config.csv")
	fmt.Println(DataConfig)

}
