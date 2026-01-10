package email

import (
	"bytes"
	"fmt"
	"html/template"
	"jogjaborobudur-chat/internal/domain/chat/entity"
)

func BuildConversationHTML(conv *entity.ChatConversation) string {
	var buf bytes.Buffer

	for _, m := range conv.Messages {
		align := "left"
		bg := "#f1f1f1"

		if m.SenderType == "admin" {
			align = "right"
			bg = "#d1e7ff"
		}

		buf.WriteString(fmt.Sprintf(`
				<div style="text-align:%s;margin:8px 0">
					<span style="display:inline-block;padding:10px;background:%s;border-radius:8px;">
						%s
					</span>
				</div>
			`, align, bg, template.HTMLEscapeString(deref(m.Message))))
	}

	return buf.String()
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safe(msg *string) string {
	if msg == nil {
		return ""
	}
	return *msg
}
