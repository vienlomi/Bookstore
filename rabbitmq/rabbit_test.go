package rabbitmq

import (
	"connection"
	"fmt"
	"testing"
	"time"
)

func TestRabbit(t *testing.T) {
	connection.Start("/home/binh/Desktop/project_v3/src/config/config.csv")
	conn, err := NewConnectMQ()
	if err != nil {
		t.Error("can not connect rabbitmq")
	}
	defer CloseMQ(conn)

	ch, err := NewChannelMQ(conn)
	if err != nil {
		t.Error("can not create channel rabbitmq")
	}

	// create new main exchange and dead-letter exchange
	err = ch.NewSender("ex1", "fanout")
	if err != nil {
		t.Error("can not create sender exchange")
	}

	//create dead letter queue for about exchange
	deadMsgs, err := ch.NewReceiverDead("deadEx", 1000)
	if err != nil {
		t.Error("can not create dead consumer")
	}
	go func() {
		fmt.Println("dead waiting")
		for d := range deadMsgs {
			fmt.Println("DEAD LETTER: ", d.Body)
			//time.Sleep(time.Second * 1)
		}
	}()

	for i := 1; i <= 10; i++ {
		msgs, err := ch.NewReceiver("ex1", "deadEx", "", 1000)
		if err != nil {
			t.Error("can not consumer msg")
		}
		go func() {
			fmt.Println("waiting")
			for d := range msgs {
				fmt.Println(d.Body)
				time.Sleep(time.Second * 5)
				d.Nack(false, false)
			}
		}()
	}

	for i := 1; i <= 100; i++ {
		err = ch.PublishMQ("ex1", "", "my name is binh ")
		if err != nil {
			t.Error("can not publish msg to exchange")
		}
		time.Sleep(time.Millisecond * 100)
	}

	time.Sleep(time.Second * 30)

}

//func TestConn(t *testing.T) {
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//	if err != nil {
//		t.Error(err, "Failed to connect to RabbitMQ")
//	}
//	fmt.Println("connect ok")
//	defer conn.Close()
//}
