package interfaces

import "jogjaborobudur-chat/internal/domain/chat/entity"

type UserChatInterface interface {
	Create(user *entity.UserChat) error
	FindByEmail(email string) (*entity.UserChat, error)
	FindBySession(session string) (*entity.UserChat, error)
	CheckExpiredUserSession(session string) (bool, error)
	DeleteExpiredUsers() error
}
