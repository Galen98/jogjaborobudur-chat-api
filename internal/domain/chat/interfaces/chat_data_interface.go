package interfaces

import "jogjaborobudur-chat/internal/domain/chat/entity"

type ChatDataInterface interface {
	SaveMessage(chat *entity.ChatData) (*entity.ChatData, error)
	GetConversationByToken(token string) (*entity.ChatConversation, error)
	GetMessagesByToken(token string, limit int, offset int) ([]entity.ChatData, error)
	//GetLastMessageByToken(token string) (*entity.ChatData, error)
	//CountMessagesByToken(token string) (int64, error)
}
