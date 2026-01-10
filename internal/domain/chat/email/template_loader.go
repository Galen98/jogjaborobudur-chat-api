package email

import (
	"os"
	"strings"
)

func LoadConversationTemplate(messagesHTML string) (string, error) {
	data, err := os.ReadFile("internal/domain/chat/email/template/conversation.html")
	if err != nil {
		return "", err
	}

	html := string(data)

	html = strings.Replace(html, "{{messages}}", messagesHTML, 1)

	return html, nil
}
