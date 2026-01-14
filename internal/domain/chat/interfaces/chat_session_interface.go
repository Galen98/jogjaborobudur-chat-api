package interfaces

import (
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/entity"
)

type ChatSessionInterface interface {
	InitChatSession(chatsession *entity.ChatSession) error
	GetChatSessionByUser(session string, productId uint) (*entity.ChatSession, error)
	GetAllChatSessionByUser(session string) ([]entity.ChatSession, error)
	GetAllChatSession() ([]dto.AdminSessionDto, error)
	UpdateSession(session *entity.ChatSession) error
	UpdateSessionOpen(session *entity.ChatSession) error
	OpenChatByUser(token string) error
	OpenChatByAdmin(token string) error
	GetChatSessionByToken(token string) (*entity.ChatSession, error)
}
