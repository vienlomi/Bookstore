package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

//consumer
type Consumer struct {
	Cchannel     *ChannelMQ
	Cname        string // name consumer get msg
	Exchange     string // exchange that we will bind to
	ExchangeType string // topic, direct, etc...
	BindingKey   string // routing key that we are using
}

func NewConsumer(ch *ChannelMQ, cname, exchange, exchangeType, bindingKey string) *Consumer {
	return &Consumer{
		Cchannel:     ch,
		Cname:        cname,
		Exchange:     exchange,
		ExchangeType: exchangeType,
		BindingKey:   bindingKey,
	}
}

func NewConsumerDead(ch *ChannelMQ) *Consumer {
	return &Consumer{
		Cchannel:     ch,
		Cname:        "deadConsumer",
		Exchange:     "deadEx",
		ExchangeType: "direct",
		BindingKey:   "deadKey",
	}
}

func (c *Consumer) Connect(cnt int) ([]<-chan amqp.Delivery, error) {
	// declare exchange
	err := c.Cchannel.Channel.ExchangeDeclare(
		c.Exchange,     // name
		c.ExchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		fmt.Println("can not declare main exchange ", err)
		return nil, err
	}
	err = c.Cchannel.Channel.ExchangeDeclare(
		"deadEx", // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		fmt.Println("can not declare dead-letter exchange ", err)
		return nil, err
	}

	//declare queue
	args := amqp.Table{
		"x-dead-letter-exchange":    "deadEx",
		"x-message-ttl":             300000,
		"x-dead-letter-routing-key": "deadKey",
	}

	q, err := c.Cchannel.Channel.QueueDeclare(
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

	err = c.Cchannel.Channel.QueueBind(
		q.Name,       // queue name
		c.BindingKey, // routing key
		c.Exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		fmt.Println("can not queue bind ", err)
		return nil, err
	}

	fmt.Println(cnt, " so chan consumers")
	var msgs []<-chan amqp.Delivery
	for i := 1; i <= cnt; i++ {
		msg, err := c.Cchannel.Channel.Consume(
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
		msgs = append(msgs, msg)
	}

	return msgs, nil

}
