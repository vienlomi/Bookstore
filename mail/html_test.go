package mail

import (
	"config"
	"connection"
	"data"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"testing"
)

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
        }
        .infor span{
            padding-right: 30px;
        }
    </style>
</head>
<body>
    <div class="order">
        <h3>Thanks for your buying our books</h3>
        <p>This is all items you have ready ordered</p>
        <div class="items">
            <ul>
                %s
            </ul>
        </div>
    </div>
</body>
</html>`, htmlItems)
}

func TestHtml(t *testing.T) {
	var order data.Order
	order.Items = append(order.Items, data.OrderItem{Name: "1", Price: 1, ImageUrl: "https://images-na.ssl-images-amazon.com/images/I/81lAPl9Fl0L.jpg", ProductId: 1, Amount: 1, UnitPrice: 1})
	order.Items = append(order.Items, data.OrderItem{Name: "2", Price: 1, ImageUrl: "https://images-na.ssl-images-amazon.com/images/I/81lAPl9Fl0L.jpg", ProductId: 1, Amount: 1, UnitPrice: 1})
	order.Items = append(order.Items, data.OrderItem{Name: "3", Price: 1, ImageUrl: "https://images-na.ssl-images-amazon.com/images/I/81lAPl9Fl0L.jpg", ProductId: 1, Amount: 1, UnitPrice: 1})
	//
	//HtmlMailOrder(order)
	connection.Start("/home/binh/Desktop/project_v3/src/config/config.csv")

	ec := EmailContent{
		Subject:      config.DataConfig["subject_mail"],
		FromUser:     &EmailUser{Name: config.DataConfig["company_name"], Email: config.DataConfig["company_mail"]},
		ToUser:       &EmailUser{Name: "user", Email: "dpbinh97@gmail.com"},
		PlainContent: "aaaaaaaaaa",
		HtmlContent:  HtmlMailOrder(order),
	}
	from := mail.NewEmail(ec.FromUser.Name, ec.FromUser.Email)
	subject := ec.Subject
	to := mail.NewEmail(ec.ToUser.Name, ec.ToUser.Email)
	plainTextContent := ec.PlainContent
	htmlContent := ec.HtmlContent
	fmt.Println(htmlContent)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
