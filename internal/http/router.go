package http

import (
	chatRoute "jogjaborobudur-chat/internal/domain/chat/route"
	"jogjaborobudur-chat/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoute(db *gorm.DB, uc *usecase.ChatUseCase) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		chat := api.Group("/chat")
		{
			chatRoute.ChatRoute(chat, db, uc)
		}
	}

	return r
}
