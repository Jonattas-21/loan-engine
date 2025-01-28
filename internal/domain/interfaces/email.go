package interfaces

import ()

type EmailSender interface {
	SendMail(subject string, content string, emailsTo ...string) error
}