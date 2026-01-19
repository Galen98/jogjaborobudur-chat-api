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
		align := "right"
		bg := "#f1f1f1"

		if m.SenderType == "admin" {
			align = "left"
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

func BuildProductCardHTML(
	productName string,
	thumbnail string,
	link string,
) string {

	return fmt.Sprintf(`
		<div style="
			border:1px solid #e5e7eb;
			border-radius:12px;
			padding:16px;
			margin-bottom:20px;
			max-width:500px;
			font-family:Arial,sans-serif;
		">
			<img src="%s" alt="%s"
				style="width:100%%;max-height:220px;object-fit:cover;border-radius:8px;" />

			<h3 style="margin:12px 0 8px 0;color:#111827;">
				%s
			</h3>

			<a href="%s"
				style="
					display:inline-block;
					padding:12px 20px;
					background:#2563eb;
					color:#ffffff;
					text-decoration:none;
					border-radius:8px;
					font-weight:bold;
				">
				Check availability
			</a>
		</div>
	`,
		thumbnail,
		template.HTMLEscapeString(productName),
		template.HTMLEscapeString(productName),
		link,
	)
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
