package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"project_v3/config"
	"project_v3/data"
	"project_v3/rabbitmq"
)

/*
Goal: Every 1 minute
	- Scan order from database -> array EmailContent
	- Add host infor, clients infor, content ... to array EmailContent
	- Put this array EmailContent into channel of sending
*/

type Prepare struct {
	channel *rabbitmq.ChannelMQ
	exName  string
	Ctx     context.Context
}

//Create new Scheduler
func NewPrepare(chann *rabbitmq.ChannelMQ, ctx context.Context, exName string, exType string) *Prepare {
	if chann == nil {
		return nil
	}
	var pre Prepare
	pre.channel = chann
	if err := pre.channel.NewSender(exName, exType); err != nil {
		fmt.Println("can not prepare new sender")
		return nil
	}
	pre.exName = exName
	pre.Ctx = ctx

	return &pre
}

//func (Scheduler) Start run
func (pre *Prepare) PreSend(username string, email string, order data.Order) {

	fmt.Println("prepare send email", pre.exName)
	var emailContent EmailContent
	emailContent.Subject = config.DataConfig["subject_mail"]
	emailContent.FromUser = &EmailUser{
		Name:  config.DataConfig["company_name"],
		Email: config.DataConfig["company_mail"],
	}
	emailContent.ToUser = &EmailUser{
		Name:  username,
		Email: email,
	}
	emailContent.PlainContent = "confirm order email"
	emailContent.HtmlContent = HtmlMailOrder(order)
	fmt.Println(emailContent)
	msgchang, err := json.Marshal(emailContent)
	if err != nil {
		fmt.Println("can't marshal email content")
	}
	err = pre.channel.PublishMQ(pre.exName, "", msgchang)
	if err != nil {
		fmt.Println("can't prepare email")
	}
}

func (pre *Prepare) PreSend2(token string, email string) {

	fmt.Println("prepare send email", pre.exName)
	var emailContent EmailContent
	emailContent.Subject = "Reset Your Password"
	emailContent.FromUser = &EmailUser{
		Name:  config.DataConfig["company_name"],
		Email: config.DataConfig["company_mail"],
	}
	emailContent.ToUser = &EmailUser{
		Name:  "Nhan Vien",
		Email: email,
	}
	emailContent.PlainContent = "Reset Password"
	emailContent.HtmlContent = HtmlMail2(token)
	fmt.Println(emailContent)
	msgchang, err := json.Marshal(emailContent)
	if err != nil {
		fmt.Println("can't marshal email content")
	}
	err = pre.channel.PublishMQ(pre.exName, "", msgchang)
	if err != nil {
		fmt.Println("can't prepare email")
	}
}

func HtmlMail2(token string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password: the 19th Bookstore</title>
</head>
<body>
    <div>
        <h3>Click the link below to reset your password.</h3>
        <a href="http://localhost:8080/reset-token/%s"> Link </a>
    </div>
</body>
</html>`, token)
}

func HtmlMailOrder(order data.Order) string {
	htmlItems := ""
	for _, item := range order.Items {
		htmlItems += fmt.Sprintf(`
		<li>
			<img src="%s">
			<div class="info">
				<p>Name: %s</p>
				<p>Amount: %d</p>
				<p>Unit Price: %f</p>
			</div>
		</li>`, item.ImageUrl, item.Name, item.Amount, item.UnitPrice)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Order confirm</title>
        <style>
        .order{
            margin-left: 20px;
        }
        img {
            height: 100px;
            margin-right: 20px;
        }
        .items ul {
            display: flex;
			flex-wrap: wrap;
        }
        .items ul li{
            display: flex;
			max-width: 30rem;
        }
        .infor span{
            padding-right: 30px;
        }
    </style>
</head>
<body>
    <div class="order">
        <h3>Thanks for your buying our books. We will delivery as soon as possible.</h3>
        <p>This is all items you have ready ordered</p>
        <div class="items">
            <ul>
                %s
            </ul>
        </div>
		<p>If you did not initiate this request, please contact us immediately at support@binhvien.com.</p>
		<p>Thank you very much</p>
    </div>
</body>
</html>`, htmlItems)
}
