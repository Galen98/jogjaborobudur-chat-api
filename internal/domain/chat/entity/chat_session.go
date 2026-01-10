package entity

import "time"

type ChatSession struct {
	ID          uint      `gorm:"primaryKey;column:id"`
	ProductID   uint      `gorm:"column:product_id;index"`
	Thumbnail   string    `gorm:"column:thumbnail"`
	ProductName string    `gorm:"column:product_name"`
	IsRead      bool      `gorm:"column:is_read"`
	IsReadAdmin bool      `gorm:"column:is_read_admin"`
	UserSession string    `gorm:"column:user_session"`
	Token       string    `gorm:"column:token"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (ChatSession) TableName() string {
	return "chat_session"
}
