package entity

import "time"

type ChatData struct {
	ID               uint      `gorm:"primaryKey;column:id"`
	ChatSessionToken string    `gorm:"column:chat_session_token"`
	Message          *string   `gorm:"column:message"`
	SenderType       string    `gorm:"column:sender_type"`
	Time             time.Time `gorm:"column:time"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (ChatData) TableName() string {
	return "chat_data"
}

type ChatConversation struct {
	Token    string     `json:"token"`
	Messages []ChatData `json:"messages"`
}
