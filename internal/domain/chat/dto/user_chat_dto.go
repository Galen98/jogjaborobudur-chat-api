package dto

import "time"

type CreateUserChatRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Session  string `json:"session"`
}

type UserChatResponse struct {
	ID          uint      `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Session     string    `json:"session"`
	ExpiredDate time.Time `json:"expired_date"`
}

type UserChatRequest struct {
	Session *string `json:"session"`
	Email   *string `json:"email"`
}
