package email

import (
	"fmt"
	"jogjaborobudur-chat/internal/domain/chat/entity"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EmailService struct {
	smtp *SMTPClient
}

func NewEmailService(smtp *SMTPClient) *EmailService {
	return &EmailService{smtp: smtp}
}

func Title(s string) string {
	return cases.Title(language.English).String(s)
}

func (s *EmailService) SendConversationEmail(
	to string,
	fullName string,
	productName string,
	conv *entity.ChatConversation,
) error {

	subject := fmt.Sprintf(
		"Chat with %s about %s",
		Title(fullName),
		Title(productName),
	)

	messagesHTML := BuildConversationHTML(conv)

	htmlBody, err := LoadConversationTemplate(messagesHTML)
	if err != nil {
		return err
	}

	return s.smtp.Send(SMTPMessage{
		To:      to,
		Subject: subject,
		HTML:    htmlBody,
	})
}
