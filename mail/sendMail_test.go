package mail

import (
	"connection"
	"context"
	"data"
	"rabbitmq"
	"testing"
	"time"
)

var pre *Prepare

func TestSend(t *testing.T) {
	connection.Start("/home/binh/Desktop/project_v3/src/config/config.csv")
	conn, err := rabbitmq.NewConnectMQ()
	if err != nil {
		t.Error("can not connect rabbitmq")
	}
	defer rabbitmq.CloseMQ(conn)

	ch, err := rabbitmq.NewChannelMQ(conn)
	if err != nil {
		t.Error("can not create channel rabbitmq")
	}

	consumerDead := rabbitmq.NewConsumerDead(ch)

	//var consumers []*rabbitmq.Consumer

	ctx := context.Background()
	mailer := NewSendgrid("apiKey")
	var order data.Order
	consumer1 := rabbitmq.NewConsumer(ch, "", "ex1", "fanout", "")
	consumer2 := rabbitmq.NewConsumer(ch, "", "ex1", "fanout", "")
	consumer3 := rabbitmq.NewConsumer(ch, "", "ex1", "fanout", "")

	//for i := 1; i <= 3; i++ {
	//	consumers = append(consumers, rabbitmq.NewConsumer(ch, "","ex1","fanout", ""))
	//}
	NewWorker(mailer, consumerDead, ctx).Start(1)
	NewWorker(mailer, consumer3, ctx).Start(1)
	NewWorker(mailer, consumer1, ctx).Start(1)
	NewWorker(mailer, consumer2, ctx).Start(1)

	//for i := range consumers {
	//	NewWorker(mailer, consumers[i], Ctx).Start()
	//}

	pre = NewPrepare(ch, ctx, "ex1", "fanout")

	pre.PreSend("binh", "dpbinh97@gmail.com", order)

	time.Sleep(time.Second * 100)
}
