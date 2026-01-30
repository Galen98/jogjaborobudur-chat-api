package repository

import (
	"errors"
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/interfaces"
	"time"

	"gorm.io/gorm"
)

var _ interfaces.ChatSessionInterface = (*ChatSessionRepository)(nil)
var ErrChatSessionNotFound = errors.New("chat session not found")

type ChatSessionRepository struct {
	db *gorm.DB
}

func NewChatSessionRepository(db *gorm.DB) interfaces.ChatSessionInterface {
	return &ChatSessionRepository{db: db}
}

func (r *ChatSessionRepository) InitChatSession(chatsession *entity.ChatSession) error {
	return r.db.Create(chatsession).Error
}

func (r *ChatSessionRepository) GetChatSessionByUser(session string, productId uint) (*entity.ChatSession, error) {
	var chatsession entity.ChatSession
	err := r.db.Where("user_session = ? AND product_id = ?", session, productId).First(&chatsession).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatSessionNotFound
		}

		return nil, err
	}
	return &chatsession, nil
}

func (r *ChatSessionRepository) UpdateSession(session *entity.ChatSession) error {
	return r.db.Model(&entity.ChatSession{}).
		Where("token = ?", session.Token).
		Updates(map[string]interface{}{
			"is_read":       session.IsRead,
			"is_read_admin": session.IsReadAdmin,
			"updated_at":    time.Now(),
		}).Error
}

func (r *ChatSessionRepository) UpdateSessionOpen(session *entity.ChatSession) error {
	return r.db.Model(&entity.ChatSession{}).
		Where("token = ?", session.Token).
		Updates(map[string]interface{}{
			"is_read":       session.IsRead,
			"is_read_admin": session.IsReadAdmin,
		}).Error
}

func (r *ChatSessionRepository) GetAllChatSessionByUser(sessionId string) ([]entity.ChatSession, error) {
	var sessions []entity.ChatSession

	err := r.db.
		Where("user_session = ?", sessionId).
		Order("updated_at DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *ChatSessionRepository) GetAllChatSession() ([]dto.AdminSessionDto, error) {
	var result []dto.AdminSessionDto

	err := r.db.Table("chat_session cs").
		Select(`
			cs.id,
			cs.link,
			cs.user_session,
			cs.product_id,
			cs.product_name,
			cs.thumbnail,
			cs.is_read,
			cs.is_read_admin,
			cs.token,
			cs.updated_at,
			uc.full_name
		`).
		Joins("LEFT JOIN user_chat uc ON uc.session = cs.user_session").
		Order("cs.updated_at DESC").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ChatDataRepository) DeleteExpiredSession() {

}

func (r *ChatSessionRepository) GetAdminSessionByToken(token string) (*dto.AdminSessionDto, error) {
	var result dto.AdminSessionDto

	err := r.db.Table("chat_session cs").
		Select(`
            cs.id,
			cs.user_session,
            cs.token,
			cs.link,
            cs.product_id,
            cs.product_name,
            cs.thumbnail,
            cs.is_read,
            cs.is_read_admin,
            cs.updated_at,
            uc.full_name
        `).
		Joins("LEFT JOIN user_chat uc ON uc.session = cs.user_session").
		Where("cs.token = ?", token).
		Scan(&result).Error

	return &result, err
}

func (r *ChatSessionRepository) GetChatSessionByToken(token string) (*entity.ChatSession, error) {
	var session entity.ChatSession

	err := r.db.
		Where("token = ?", token).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *ChatSessionRepository) OpenChatByUser(token string) error {
	return r.db.
		Model(&entity.ChatSession{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_read": true,
		}).Error
}

func (r *ChatSessionRepository) OpenChatByAdmin(token string) error {
	return r.db.
		Model(&entity.ChatSession{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_read_admin": true,
		}).Error
}
