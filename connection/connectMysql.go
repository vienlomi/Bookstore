package connection

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func StartConnectMysql() {
	db, err := sql.Open("mysql", "vien:vienlomi@/project_v3")
	if err != nil {
		fmt.Println("Can not connect mysql")
	}

	DB = db
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(100)
	DB.SetConnMaxLifetime(time.Second * 5)
	DB.SetConnMaxIdleTime(time.Second * 5)
}

func Disconnect() error {
	err := DB.Close()
	return err
}
