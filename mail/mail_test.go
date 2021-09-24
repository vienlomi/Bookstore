package mail

import (
	"connection"
	"data"
	"testing"
	"time"
)

func TestCon(t *testing.T) {
	connection.Start("/home/binh/Desktop/project_v3/src/config/config.csv")
	StartMQ()
	var order data.Order

	for i := 1; i <= 10; i++ {
		Pre.PreSend("a", "b", order)
		time.Sleep(time.Second * 1)

	}
	time.Sleep(time.Second * 20)

}

//func TestSendGmail(t *testing.T) {
//	from := mail.NewEmail("Binh", "binh92dangdung@gmail.com") // Change to your verified sender
//	subject := "Sending with Twilio SendGrid is Fun"
//	to := mail.NewEmail("dpbinh", "dpbinh97@gmail.com") // Change to your recipient
//	plainTextContent := "and easy to do anywhere, even with Go"
//	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
//	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
//	client := sendgrid.NewSendClient(os.Getenv("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"))
//
//	response, err := client.Send(message)
//	if err != nil {
//		log.Println(err)
//	} else {
//		fmt.Println(response.StatusCode)
//		fmt.Println(response.Headers)
//	}
//}
