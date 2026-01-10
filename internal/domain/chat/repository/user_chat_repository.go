package repository

import (
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/interfaces"
	"time"

	"gorm.io/gorm"
)

type UserChatRepository struct {
	db *gorm.DB
}

func NewUserChatRepository(db *gorm.DB) interfaces.UserChatInterface {
	return &UserChatRepository{db: db}
}

func (r *UserChatRepository) Create(user *entity.UserChat) error {
	return r.db.Create(user).Error
}

func (r *UserChatRepository) FindByEmail(email string) (*entity.UserChat, error) {
	var user entity.UserChat
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserChatRepository) FindBySession(session string) (*entity.UserChat, error) {
	var user entity.UserChat
	err := r.db.Where("session = ?", session).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserChatRepository) CheckExpiredUserSession(session string) (bool, error) {
	var user entity.UserChat

	err := r.db.Where("session = ?", session).First(&user).Error

	if err != nil {
		return false, err
	}

	today := time.Now().Truncate(24 * time.Hour)

	if user.ExpiredDate.Before(today) {
		return true, nil
	}

	return false, nil
}
