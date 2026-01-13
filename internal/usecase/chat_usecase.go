package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/email"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/services"
	"jogjaborobudur-chat/internal/ws"
	"time"

	"github.com/redis/go-redis/v9"
)

type ChatUseCase struct {
	chatDataService    *services.ChatDataService
	chatSessionService *services.ChatSessionService
	userChatService    *services.UserChatService
	emailService       *email.EmailService
	Hub                *ws.Hub
	redis              *redis.Client
}

func NewChatUseCase(
	chatDataService *services.ChatDataService,
	chatSessionService *services.ChatSessionService,
	userChatService *services.UserChatService,
	emailService *email.EmailService,
	hub *ws.Hub,
	redis *redis.Client,
) *ChatUseCase {
	return &ChatUseCase{
		chatDataService:    chatDataService,
		chatSessionService: chatSessionService,
		userChatService:    userChatService,
		emailService:       emailService,
		Hub:                hub,
		redis:              redis,
	}
}

func (u *ChatUseCase) SendMessage(
	req dto.SendChatRequest,
) (*entity.ChatData, error) {

	session, err := u.chatSessionService.GetByToken(req.Token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("token required")
	}

	// simpan message
	msg, err := u.chatDataService.SendMessage(req)
	if err != nil {
		return nil, err
	}

	// update read status
	if req.SenderType == "user" {
		session.IsRead = true
		session.IsReadAdmin = false
	} else {
		session.IsRead = false
		session.IsReadAdmin = true
	}

	session.UpdatedAt = time.Now()

	if err := u.chatSessionService.UpdateSessionStatus(session); err != nil {
		return nil, err
	}

	// 🔥 PUBLISH KE REDIS (BUKAN WS)
	msgPayload, _ := json.Marshal(msg)
	sessionPayload, _ := json.Marshal(session)

	ctx := context.Background()

	// realtime chat room
	u.redis.Publish(ctx, "chat:"+session.Token, msgPayload)

	// realtime user chat list
	u.redis.Publish(ctx, "session:"+session.UserSession, sessionPayload)

	// realtime admin dashboard
	u.redis.Publish(ctx, "admin:sessions", sessionPayload)

	// email (side effect, OK)
	if req.SenderType == "admin" {
		user, err := u.userChatService.GetBySession(session.UserSession)
		if err == nil && user.Email != "" {
			conv, _ := u.chatDataService.GetConversationAll(req.Token)
			go u.emailService.SendConversationEmail(
				user.Email,
				user.FullName,
				session.ProductName,
				conv,
			)
		}
	}

	return msg, nil
}

func (u *ChatUseCase) GetMessagesByToken(token string, limit int, offset int) ([]entity.ChatData, error) {
	if token == "" {
		return nil, errors.New("token required")
	}

	return u.chatDataService.GetMessagesPaginated(token, limit, offset)
}

func (u *ChatUseCase) OpenChatByUser(token string, types string) error {
	session, err := u.chatSessionService.GetByToken(token)

	if err != nil {
		return err
	}

	if types == "user" {
		session.IsRead = true
	} else {
		session.IsReadAdmin = true
	}

	if err := u.chatSessionService.UpdateSessionStatus(session); err != nil {
		return err
	}

	payload, _ := json.Marshal(session)

	u.Hub.Broadcast(ws.BroadcastMessage{
		Token: "session:" + session.UserSession,
		Data:  payload,
	})

	u.Hub.Broadcast(ws.BroadcastMessage{
		Token: "admin:sessions",
		Data:  payload,
	})

	return nil
}

func (u *ChatUseCase) GetUserSession(
	sessionId string,
	productId uint,
) (*entity.ChatSession, error) {
	return u.chatSessionService.GetChatSessionByUser(sessionId, productId)
}

func (u *ChatUseCase) GetAdminSessions() ([]dto.AdminSessionDto, error) {
	return u.chatSessionService.GetAllChatSession()
}

func (u *ChatUseCase) GetUserChatSessionByUser(sessionId string) ([]entity.ChatSession, error) {
	return u.chatSessionService.GetAllChatSessionByUser(sessionId)
}
