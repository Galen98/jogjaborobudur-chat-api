package usecase

import (
	"encoding/json"
	"errors"
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/email"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/services"
	"jogjaborobudur-chat/internal/ws"
	"time"
)

type ChatUseCase struct {
	chatDataService    *services.ChatDataService
	chatSessionService *services.ChatSessionService
	userChatService    *services.UserChatService
	emailService       *email.EmailService
	Hub                *ws.Hub
}

func NewChatUseCase(
	chatDataService *services.ChatDataService,
	chatSessionService *services.ChatSessionService,
	userChatService *services.UserChatService,
	emailService *email.EmailService,
	hub *ws.Hub,
) *ChatUseCase {
	return &ChatUseCase{
		chatDataService:    chatDataService,
		chatSessionService: chatSessionService,
		userChatService:    userChatService,
		emailService:       emailService,
		Hub:                hub,
	}
}

func (u *ChatUseCase) SendMessage(req dto.SendChatRequest) (*entity.ChatData, error) {

	session, err := u.chatSessionService.GetByToken(req.Token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("token required")
	}

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

	msgPayload, _ := json.Marshal(msg)

	u.Hub.Broadcast(ws.BroadcastMessage{
		Token: "chat:" + session.Token,
		Data:  msgPayload,
	})

	sessionPayload, _ := json.Marshal(session)

	u.Hub.Broadcast(ws.BroadcastMessage{
		Token: "session:" + session.UserSession,
		Data:  sessionPayload,
	})

	u.Hub.Broadcast(ws.BroadcastMessage{
		Token: "admin:sessions",
		Data:  sessionPayload,
	})

	_ = u.chatSessionService.UpdateSessionStatus(session)

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
