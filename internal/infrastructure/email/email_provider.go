package email

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
	"log"
)

type EmailSender struct {
}

func (e *EmailSender) SendMail(subject string, content string, emailsTo ...string) error {

	d := gomail.NewDialer(os.Getenv("MAIL_SMTP_HOST"), 587, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASSWORD"))
	log.Println("Sending email to: ", emailsTo)

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("MAIL_USER"))
	message.SetHeader("To", emailsTo...)
	message.SetHeader("Subject", fmt.Sprint("Hello! ", subject))
	message.SetBody("text/html", content)

	err := d.DialAndSend(message)
	log.Println("Email sent to: ", emailsTo)
	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}

	return nil
}
