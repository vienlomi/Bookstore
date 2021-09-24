package rabbitmq

import (
	"connection"
	"fmt"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
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
	consumerDead := NewConsumerDead(ch)
	msgDead, err := consumerDead.Connect(1)
	if err != nil {
		t.Error("consumerDead fail")
	}
	go func() {
		fmt.Println("dead consumer waiting ....")
		for d := range msgDead {
			fmt.Println("DEAD LETTER: ", d.Body)
		}
	}()

	consumer1 := NewConsumer(ch, "consumer1", "ex1", "fanout", "")
	msgs, err := consumer1.Connect()
	if err != nil {
		t.Error("consumer1 fail")
	}
	go func() {
		fmt.Println("consumer waiting ....")
		for d := range msgs {
			fmt.Println(d.Body)
		}
		fmt.Println("Done :", consumer1.Done)
	}()

	for i := 1; i <= 10; i++ {
		err = ch.PublishMQ("ex1", "", "my name is binh")
		if err != nil {
			t.Error("can not publish msg to exchange")
		}
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(time.Second * 20)
}
