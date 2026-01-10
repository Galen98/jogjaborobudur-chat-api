package dto

import "time"

type GetChatHistoryRequest struct {
	Token string `json:"token" binding:"required"`
}

type ChatHistoryItem struct {
	ID         uint      `json:"id"`
	Message    string    `json:"message"`
	SenderType string    `json:"sender_type"`
	Time       time.Time `json:"time"`
}

type GetChatHistoryResponse struct {
	Token    string            `json:"token"`
	Messages []ChatHistoryItem `json:"messages"`
}

type SendChatRequest struct {
	Token      string    `json:"token" binding:"required"`
	Message    string    `json:"message" binding:"required"`
	SenderType string    `json:"sender_type" binding:"required"`
	Time       time.Time `json:"time"`
}

type ChatStreamEvent struct {
	Event      string    `json:"event"`
	Token      string    `json:"token"`
	Message    string    `json:"message"`
	SenderType string    `json:"sender_type"`
	Time       time.Time `json:"time"`
}

type GetChatHistoryPaginatedRequest struct {
	Token  string `json:"token"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
