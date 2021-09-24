package mail

import (
	"config"
	"connection"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"testing"
)

func TestSenddddd(t *testing.T) {
	connection.Start("/home/binh/Desktop/project_v3/src/config/config.csv")

	ec := EmailContent{
		Subject:      config.DataConfig["subject_mail"],
		FromUser:     &EmailUser{Name: config.DataConfig["company_name"], Email: config.DataConfig["company_mail"]},
		ToUser:       &EmailUser{Name: "user", Email: "dpbinh97@gmail.com"},
		PlainContent: "aaaaaaaaaa",
		HtmlContent:  "bbbbbbbb",
	}
	from := mail.NewEmail(ec.FromUser.Name, ec.FromUser.Email)
	subject := ec.Subject
	to := mail.NewEmail(ec.ToUser.Name, ec.ToUser.Email)
	plainTextContent := ec.PlainContent
	htmlContent := ec.HtmlContent
	fmt.Println(subject, from, to)
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
