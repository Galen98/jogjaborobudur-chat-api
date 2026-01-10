package repository

import (
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/interfaces"

	"gorm.io/gorm"
)

type ChatDataRepository struct {
	db *gorm.DB
}

func NewChatDataRepository(db *gorm.DB) interfaces.ChatDataInterface {
	return &ChatDataRepository{db: db}
}

func (r *ChatDataRepository) SaveMessage(chat *entity.ChatData) (*entity.ChatData, error) {
	err := r.db.Create(chat).Error
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *ChatDataRepository) GetConversationByToken(token string) (*entity.ChatConversation, error) {
	var messages []entity.ChatData

	err := r.db.
		Where("chat_session_token = ?", token).
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	conv := &entity.ChatConversation{
		Token:    token,
		Messages: messages,
	}

	return conv, nil
}

func (r *ChatDataRepository) GetMessagesByToken(token string, limit int, offset int) ([]entity.ChatData, error) {
	var messages []entity.ChatData

	err := r.db.Where("chat_session_token = ?", token).Order("created_at DESC").Limit(limit).Offset(offset).Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
