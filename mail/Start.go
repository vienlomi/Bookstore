package mail

import (
	"context"
	"fmt"
	"project_v3/config"
	"project_v3/connection"
	"project_v3/rabbitmq"
	"strconv"

	"github.com/streadway/amqp"
)

var (
	Conn    *amqp.Connection
	Channel *rabbitmq.ChannelMQ
	C       *rabbitmq.Consumer
	Pre     *Prepare
)

func StartMQ() {

	fmt.Println("StartMQ : ", Conn, Channel, C, Pre)
	connection.Start("/home/binh/Desktop/project_v3/src/config/config.csv")
	Con, err := rabbitmq.NewConnectMQ()
	if err != nil {
		fmt.Println("can not connect rabbitmq")
	}
	//defer rabbitmq.CloseMQ(Con)

	Channel, err = rabbitmq.NewChannelMQ(Con)
	if err != nil {
		fmt.Println("can not create channel rabbitmq")
	}
	fmt.Println("StartMQ ch...", Channel)

	mailer := NewSendgrid(config.DataConfig["api_key"])
	ctx, _ := context.WithCancel(context.Background())
	//cntC := config.DataConfig["cntConsumer"]
	//cnt, _ := strconv.Atoi(cntC)

	consumerDead := rabbitmq.NewConsumerDead(Channel)
	NewWorker(mailer, consumerDead, ctx).Start(1)

	cntConf := config.DataConfig["cnt_consumers"]
	cnt, _ := strconv.Atoi(cntConf)

	C = rabbitmq.NewConsumer(Channel, "", "ex1", "fanout", "")
	NewWorker(mailer, C, ctx).Start(cnt)

	Pre = NewPrepare(Channel, ctx, "ex1", "fanout")
	fmt.Println("StartMQ pre...", Pre)

}

func CloseMQ() {
	Conn.Close()
}
