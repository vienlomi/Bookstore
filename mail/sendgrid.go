package mail

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

/*
Goal:

*/

type Emailer interface {
	Send(*EmailContent) error
}

//MailUser define email infor
type EmailUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

//change MailUser to string of json
func (eu *EmailUser) String() string {
	c, _ := json.Marshal(eu)
	return string(c)
}

//EmailContent define email content infor
type EmailContent struct {
	Subject      string     `json:"subject"`       //subject of email
	FromUser     *EmailUser `json:"from_user"`     //Email of my company
	ToUser       *EmailUser `json:"to_user"`       //Email of client
	PlainContent string     `json:"plain_content"` //Body content of email
	HtmlContent  string     `json:"html_content"`  //Html show
}

//func change EmailContent to string of json => for sending into MQ
func (ec *EmailContent) String() string {
	c, _ := json.Marshal(ec)
	return string(c)
}

//func Validate will check email content is available or not
func (ec *EmailContent) Validate() error {
	if ec == nil || ec.FromUser == nil || ec.ToUser == nil || ec.PlainContent == "" {
		return errors.New("wrong content email")
	}
	return nil
}

//Struct Sendgrid implement send email to destination email via method Send
type Sendgrid struct {
	ApiKey string `json:"api_key"`
	Client *sendgrid.Client
}

//fun NewSendgrid create new Sendgrid
func NewSendgrid(apiKey string) *Sendgrid {
	client := sendgrid.NewSendClient(apiKey)
	return &Sendgrid{
		ApiKey: apiKey,
		Client: client,
	}
}

//func Send will send to client's email base on email content
func (sg *Sendgrid) Send(ec *EmailContent) error {

	fmt.Println(ec)
	if err := ec.Validate(); err != nil {
		fmt.Println("Error check validate when sending")
	}

	from := mail.NewEmail(ec.FromUser.Name, "nhanhongvien199@gmail.com")
	subject := ec.Subject
	to := mail.NewEmail(ec.ToUser.Name, ec.ToUser.Email)
	plainTextContent := ec.PlainContent
	htmlContent := ec.HtmlContent
	fmt.Println("chuan bi gui mail: ", htmlContent)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient("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")

	response, err := client.Send(message)
	if err != nil {
		log.Println("cos loi ", err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}
