package email

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Logger *logrus.Logger
}

func (e *EmailSender) SendMail(subject string, content string, emailsTo ...string) error {

	d := gomail.NewDialer(os.Getenv("MAIL_SMTP"), 587, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASSWORD"))

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("MAIL_USER"))
	message.SetHeader("To", emailsTo...)
	message.SetHeader("Subject", fmt.Sprint("Hello! ", subject))
	message.SetBody("text/html", content)

	err := d.DialAndSend(message)
	if err != nil {
		e.Logger.Error("Error sending email: ", err)
		return err
	}

	return nil

}
