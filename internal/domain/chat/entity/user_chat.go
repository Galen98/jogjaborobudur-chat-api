package entity

import "time"

type UserChat struct {
	ID          uint      `gorm:"primaryKey;column:id"`
	FullName    string    `gorm:"column:full_name"`
	Email       string    `gorm:"column:email"`
	Session     string    `gorm:"column:session"`
	ExpiredDate time.Time `gorm:"column:expired_date"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (UserChat) TableName() string {
	return "user_chat"
}
