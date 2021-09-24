package rabbitmq

import (
	"project_v3/config"
	//"config"
	"fmt"

	"github.com/streadway/amqp"
)

func NewConnectMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial(config.DataConfig["addressMQ"])
	if err != nil || conn == nil {
		fmt.Println("can't connect rabbitmq")
		return nil, err
	}
	return conn, nil
}
func CloseMQ(conn *amqp.Connection) {
	conn.Close()
}
