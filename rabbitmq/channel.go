package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

// define Channel MQ
type ChannelMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Ctx     context.Context
}

//create new channelMQ
func NewChannelMQ(conn *amqp.Connection) (*ChannelMQ, error) {
	if conn == nil {
		return nil, errors.New("nil connection")
	}
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Can't open channel rabbitMQ")
		return nil, err
	}
	ctx := context.Background()
	return &ChannelMQ{
		Conn:    conn,
		Channel: ch,
		Ctx:     ctx,
	}, nil
}

//create exchange fanout send msg to (queue -> consumer)
func (ch *ChannelMQ) NewSender(nameEx string, kind string) error {

	err := ch.Channel.ExchangeDeclare(
		nameEx, // name
		kind,   // type
		true,   // durable
		false,  // auto-deleted
		false,  // internal
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		fmt.Println("can not declare main exchange ", err)
		return err
	}
	return nil
}

//publish msg to exchange of channel
func (ch *ChannelMQ) PublishMQ(nameEx string, key string, msg []byte) error {
	fmt.Println("publish MQ")

	err := ch.Channel.Publish(
		nameEx, // exchange
		key,    // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		fmt.Println("can not publish to exchange ", err)
		return err
	}
	fmt.Println("publish MQ oke")
	return nil
}

//create receiver fanout (binding, queue, consumer) -> output: chan receive msg
func (ch *ChannelMQ) NewReceiver(nameEx string, nameDeadEx string, key string, ttl time.Duration) (<-chan amqp.Delivery, error) {
	args := make(map[string]interface{})
	args["x-dead-letter-exchange"] = nameDeadEx
	args["x-message-ttl"] = ttl
	args["x-dead-letter-routing-key"] = key

	q, err := ch.Channel.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		args,  // arguments
	)
	if err != nil {
		fmt.Println("can not declare queue ", err)
		return nil, err
	}

	err = ch.Channel.QueueBind(
		q.Name, // queue name
		"",     // routing key
		nameEx, // exchange
		false,
		nil,
	)
	if err != nil {
		fmt.Println("can not queue bind ", err)
		return nil, err
	}
	msgs, err := ch.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Println("can not delivery ", err)
		return nil, err
	}

	return msgs, nil
}

func (ch *ChannelMQ) NewReceiverDead(nameEx string, ttl time.Duration) (<-chan amqp.Delivery, error) {
	args := make(map[string]interface{})
	args["x-message-ttl"] = ttl
	q, err := ch.Channel.QueueDeclare(
		nameEx,
		true,
		false,
		false,
		false,
		args)
	if err != nil {
		fmt.Println("can not create queue dead letter 2222222222", err)
		return nil, err
	}

	err = ch.Channel.QueueBind(
		q.Name,    // queue name
		"key_xdl", // routing key
		"exDead1", // exchange
		false,
		nil,
	)
	if err != nil {
		fmt.Println("can not queue bind ", err)
		return nil, err
	}

	msgs, err := ch.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Println("can not dead msg ", err)
		return nil, err
	}
	return msgs, nil
}
