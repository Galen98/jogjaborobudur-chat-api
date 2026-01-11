package services

import (
	"errors"
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/interfaces"
	"time"
)

type UserChatService struct {
	repo interfaces.UserChatInterface
}

func NewChatService(repo interfaces.UserChatInterface) *UserChatService {
	return &UserChatService{repo: repo}
}

func (s *UserChatService) CreateUser(req dto.CreateUserChatRequest) (*dto.UserChatResponse, error) {
	if req.Email == "" || req.FullName == "" || req.Session == "" {
		return nil, errors.New("email, fullname wajib diisi")
	}
	expiredAt := time.Now().AddDate(0, 0, 5)
	user := &entity.UserChat{
		FullName:    req.FullName,
		Email:       req.Email,
		Session:     req.Session,
		ExpiredDate: expiredAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	resp := &dto.UserChatResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		Session:     user.Session,
		ExpiredDate: user.ExpiredDate,
	}

	return resp, nil
}

func (s *UserChatService) GetBySession(session string) (*dto.UserChatResponse, error) {
	if session == "" {
		return nil, errors.New("session required")
	}

	user, err := s.repo.FindBySession(session)
	if err != nil {
		return nil, err
	}

	resp := &dto.UserChatResponse{
		ID:          user.ID,
		Session:     user.Session,
		FullName:    user.FullName,
		Email:       user.Email,
		ExpiredDate: user.ExpiredDate,
	}

	return resp, nil
}

func (s *UserChatService) GetByEmail(email string) (*dto.UserChatResponse, error) {
	if email == "" {
		return nil, errors.New("email required")
	}

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	resp := &dto.UserChatResponse{
		FullName:    user.FullName,
		Email:       user.Email,
		ExpiredDate: user.ExpiredDate,
	}

	return resp, nil
}

func (s *UserChatService) CheckExpired(session string) (bool, error) {
	if session == "" {
		return false, errors.New("session required")
	}

	return s.repo.CheckExpiredUserSession(session)
}
