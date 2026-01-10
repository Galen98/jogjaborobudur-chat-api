package email

import (
	"fmt"
	"net/smtp"
)

type SMTPClient struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

type SMTPMessage struct {
	To      string
	Subject string
	HTML    string
}

func NewSMTPClient(host string, port int, user, pass, from string) *SMTPClient {
	return &SMTPClient{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
		From: from,
	}
}

func (c *SMTPClient) Send(msg SMTPMessage) error {
	auth := smtp.PlainAuth("", c.User, c.Pass, c.Host)

	header := ""
	header += fmt.Sprintf("From: %s\r\n", c.From)
	header += fmt.Sprintf("To: %s\r\n", msg.To)
	header += fmt.Sprintf("Subject: %s\r\n", msg.Subject)
	header += "MIME-Version: 1.0\r\n"
	header += "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"

	body := header + msg.HTML

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", c.Host, c.Port),
		auth,
		c.User,
		[]string{msg.To},
		[]byte(body),
	)
}
