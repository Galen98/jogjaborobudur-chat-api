package route

import (
	"jogjaborobudur-chat/internal/domain/chat/controller"
	"jogjaborobudur-chat/internal/domain/chat/repository"
	"jogjaborobudur-chat/internal/domain/chat/services"
	"jogjaborobudur-chat/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ChatRoute(r *gin.RouterGroup, db *gorm.DB, uc *usecase.ChatUseCase) {
	userRepo := repository.NewUserChatRepository(db)
	userService := services.NewChatService(userRepo)
	userCtrl := controller.NewUserChatController(userService)

	chatSessionRepo := repository.NewChatSessionRepository(db)
	chatSessionService := services.NewChatSessionService(chatSessionRepo)
	chatSessionCtrl := controller.NewChatSessionController(chatSessionService)

	useChatCtrl := controller.NewChatUseCaseController(uc)

	r.GET("/delete-expired-users", userCtrl.DeleteExpiredUsersSheduler)
	r.POST("/user-chat", userCtrl.CreateUserChat)
	r.POST("/user-chat-expired", userCtrl.CheckExpiredSession)
	r.GET("/user-chat", userCtrl.GetUserBySession)
	r.GET("/user-chat-email", userCtrl.GetUserByEmail)

	r.POST("/init-session", chatSessionCtrl.InitChatSession)
	r.GET("/user-session", useChatCtrl.GetUserSession)

	r.GET("/admin/chat-sessions", useChatCtrl.GetAdminSessions)
	r.GET("/user/chat-session", useChatCtrl.GetAllUserChatSession)
	r.POST("/send", useChatCtrl.SendMessage)
	r.GET("/messages", useChatCtrl.GetMessages)
	r.PATCH("/open/user/:token/:types", useChatCtrl.OpenChatByUser)

}
