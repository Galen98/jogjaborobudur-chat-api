package services

import (
	"context"
	"encoding/json"
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/interfaces"
	"jogjaborobudur-chat/internal/infrastructure/cache"
	"jogjaborobudur-chat/internal/ws"
	"time"

	"github.com/redis/go-redis/v9"
)

type ChatDataService struct {
	sessionRepo interfaces.ChatSessionInterface
	messageRepo interfaces.ChatDataInterface
	wsHub       *ws.Hub
	cache       *cache.ChatMessageCache
	redis       *redis.Client
}

func NewChatDataService(
	sessionRepo interfaces.ChatSessionInterface,
	messageRepo interfaces.ChatDataInterface,
	wsHub *ws.Hub,
	cache *cache.ChatMessageCache,
	redis *redis.Client,
) *ChatDataService {
	return &ChatDataService{
		sessionRepo: sessionRepo,
		messageRepo: messageRepo,
		wsHub:       wsHub,
		cache:       cache,
		redis:       redis,
	}
}

func (s *ChatDataService) SendMessage(
	req dto.SendChatRequest,
) (*entity.ChatData, error) {

	chat := &entity.ChatData{
		ChatSessionToken: req.Token,
		Message:          &req.Message,
		SenderType:       req.SenderType,
		Time:             time.Now(),
	}

	saved, err := s.messageRepo.SaveMessage(chat)
	if err != nil {
		return nil, err
	}

	_ = s.cache.PushMessage(req.Token, *saved)

	payload, _ := json.Marshal(saved)
	s.redis.Publish(
		context.Background(),
		"chat:"+req.Token,
		payload,
	)

	return saved, nil
}

func (s *ChatDataService) GetConversation(token string) (*entity.ChatConversation, error) {
	if conv, err := s.cache.Get(token); err == nil && conv != nil {
		return conv, nil
	}
	conv, err := s.messageRepo.GetConversationByToken(token)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(conv)

	return conv, nil
}

func (s *ChatDataService) GetConversationAll(token string) (*entity.ChatConversation, error) {
	conv, err := s.messageRepo.GetConversationByToken(token)
	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (s *ChatDataService) GetMessagesPaginated(
	token string,
	limit int,
	offset int,
) ([]entity.ChatData, error) {

	return s.messageRepo.GetMessagesByToken(token, limit, offset)
}
