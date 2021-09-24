package config

import (
	"encoding/csv"
	"fmt"
	"os"
)

var DataConfig = make(map[string] string)

func ReadConfigCsv(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("open file csv fail------",err)
		return
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Println("read file csv fail-----")
		return
	}
	for _, line := range lines {
		DataConfig[line[0]] = line[1]
	}
	return
}
