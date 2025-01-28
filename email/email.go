package email

import (
	"net/smtp"
)

type EmailSender interface {
	SendMail(to, message string) error
}

type emailSender struct {
	from string
	pw   string
	addr string
	auth smtp.Auth
}

func NewMailSender(from, pw, host, port string) EmailSender {
	auth := smtp.PlainAuth("", from, pw, host)
	addr := host + ":" + port
	return &emailSender{from, pw, addr, auth}
}

func (e *emailSender) SendMail(to string, message string) error {
	err := smtp.SendMail(e.addr, e.auth, e.from, []string{to}, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
