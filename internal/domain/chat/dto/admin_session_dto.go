package dto

import "time"

type AdminSessionDto struct {
	ID          uint      `json:"id"`
	UserSession string    `json:"user_session"`
	FullName    string    `json:"fullname"`
	ProductName string    `json:"product_name"`
	ProductID   string    `json:"product_id"`
	Thumbnail   string    `json:"thumbnail"`
	IsRead      bool      `json:"is_read"`
	IsReadAdmin bool      `json:"is_read_admin"`
	Token       string    `json:"token"`
	UpdatedAt   time.Time `json:"updated_at"`
}
