package usecase

import (
	"errors"
	"jogjaborobudur-chat/config"
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/email"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/domain/chat/services"
	"time"
)

type ChatUseCase struct {
	chatDataService    *services.ChatDataService
	chatSessionService *services.ChatSessionService
	userChatService    *services.UserChatService
	emailService       *email.EmailService
	pushService        *services.AdminPushService
}

func NewChatUseCase(
	chatDataService *services.ChatDataService,
	chatSessionService *services.ChatSessionService,
	userChatService *services.UserChatService,
	emailService *email.EmailService,
	pushService *services.AdminPushService,
) *ChatUseCase {
	return &ChatUseCase{
		chatDataService:    chatDataService,
		chatSessionService: chatSessionService,
		userChatService:    userChatService,
		emailService:       emailService,
		pushService:        pushService,
	}
}

func (u *ChatUseCase) pushNotifyAdminIfNeeded(
	session *entity.ChatSession,
	msg *entity.ChatData,
	req dto.SendChatRequest,
) {
	if req.SenderType != "user" {
		return
	}

	if u.pushService == nil {
		return
	}

	if session.IsReadAdmin {
		return
	}

	go func() {
		_ = u.pushService.NotifyNewChat(
			session.ProductName,
			*msg.Message,
		)
	}()
}

func (u *ChatUseCase) SendMessage(req dto.SendChatRequest) (*entity.ChatData, error) {

	session, err := u.chatSessionService.GetByToken(req.Token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("token required")
	}

	// simpan pesan
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

	adminSession, err := u.chatSessionService.GetAdminSessionByToken(session.Token)
	if err != nil {
		return nil, err
	}

	_ = config.Pusher.Trigger(
		"chat-"+session.Token,
		"new-message",
		msg,
	)

	_ = config.Pusher.Trigger(
		"session-"+session.UserSession,
		"session-update",
		session,
	)
	_ = config.Pusher.Trigger(
		"admin-sessions",
		"session-update",
		adminSession,
	)

	u.pushNotifyAdminIfNeeded(session, msg, req)

	if req.SenderType == "admin" {
		user, err := u.userChatService.GetBySession(session.UserSession)
		if err == nil && user.Email != "" {
			conv, _ := u.chatDataService.GetConversationAll(req.Token)
			go u.emailService.SendConversationEmail(
				user.Email,
				user.FullName,
				session.ProductName,
				session.Thumbnail,
				session.Link,
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

	if err := u.chatSessionService.UpdateSessionStatusOpen(session); err != nil {
		return err
	}

	// adminSession, err := u.chatSessionService.GetAdminSessionByToken(session.Token)
	// if err != nil {
	// 	return nil
	// }
	// _ = config.Pusher.Trigger(
	// 	"session-"+session.UserSession,
	// 	"session-update",
	// 	session,
	// )

	// _ = config.Pusher.Trigger(
	// 	"admin-sessions",
	// 	"session-update",
	// 	adminSession,
	// )

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
