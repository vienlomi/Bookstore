package connection

import (
	"fmt"
	"project_v3/config"
	"project_v3/data"
)

var Model *data.Models

func Start(path string) {
	config.ReadConfigCsv(path)
	fmt.Println(config.DataConfig)
	StartConnectMysql()
	Model = data.NewModels(DB)
	fmt.Println("start ", Model)
}
