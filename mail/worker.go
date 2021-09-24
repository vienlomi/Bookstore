package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"project_v3/rabbitmq"
)

//struct Worker define a worker
type Worker struct {
	emailer  *Sendgrid
	consumer *rabbitmq.Consumer
}

//func NewWorker create a new Worker
func NewWorker(emailer *Sendgrid, consumer *rabbitmq.Consumer, ctx context.Context) *Worker {
	return &Worker{
		emailer:  emailer,
		consumer: consumer,
	}
}

/*
func (Worker) Start starts Worker
process logic:
	1. Wait message in channel Delivery
	2. Send email with each Emailer
*/
func (w *Worker) Start(cnt int) {

	if w.emailer == nil || w.consumer == nil {
		fmt.Println("Can't start Worker with emailer nil (or consumer nil)")
		return
	}
	msgs, err := w.consumer.Connect(cnt)
	if err != nil {
		fmt.Println("create consumerDead fail")
		return
	}
	fmt.Println("workinggggggggggggg", len(msgs))
	for _, msg := range msgs {
		go func() {
			fmt.Println(" waiting ....")
			for d := range msg {
				fmt.Println("co mes trong queue", d.Body)
				var ec EmailContent

				err := json.Unmarshal(d.Body, &ec)

				fmt.Println("unmarshal : ", ec)
				if err != nil {
					fmt.Println("Can't unMarshal body msg")
					// nem di
					continue
				}
				fmt.Println("marshal: ", ec)
				err = w.emailer.Send(&ec)
				if err != nil {
					fmt.Println("Can't send msg in sendgrid")
					// lam lai
				}
			}
			fmt.Println("Channel is closed, prepare reconnect...")

		}()

	}
}
