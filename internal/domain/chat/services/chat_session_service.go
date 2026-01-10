package services

import (
	"errors"
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/interfaces"
)

type ChatSessionService struct {
	repo interfaces.ChatSessionInterface
}

func NewChatSessionService(repo interfaces.ChatSessionInterface) *ChatSessionService {
	return &ChatSessionService{repo: repo}
}

func (s *ChatSessionService) GetChatSessionByUser(
	sessionId string,
	productId uint,
) (*entity.ChatSession, error) {

	return s.repo.GetChatSessionByUser(sessionId, productId)
}

func (s *ChatSessionService) InitChatSession(req dto.CreateChatSessionRequest) error {
	session := &entity.ChatSession{
		Token:       req.Token,
		ProductID:   req.ProductId,
		ProductName: req.ProductName,
		UserSession: req.Session,
		Thumbnail:   req.Thumbnail,
		IsRead:      true,
		IsReadAdmin: false,
	}

	return s.repo.InitChatSession(session)
}

func (s *ChatSessionService) GetByToken(token string) (*entity.ChatSession, error) {

	if token == "" {
		return nil, errors.New("token is required")
	}

	session, err := s.repo.GetChatSessionByToken(token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *ChatSessionService) UpdateSessionStatus(session *entity.ChatSession) error {
	return s.repo.UpdateSession(session)
}

func (s *ChatSessionService) GetAllChatSession() ([]dto.AdminSessionDto, error) {
	return s.repo.GetAllChatSession()
}

func (s *ChatSessionService) GetAllChatSessionByUser(sessionId string) ([]entity.ChatSession, error) {
	return s.repo.GetAllChatSessionByUser(sessionId)
}

func (s *ChatSessionService) OpenChatByUser(token string) error {
	return s.repo.OpenChatByUser(token)
}

func (s *ChatSessionService) OpenChatByAdmin(token string) error {
	return s.repo.OpenChatByAdmin(token)
}
