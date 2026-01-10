package controller

import (
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/repository"
	"jogjaborobudur-chat/internal/domain/chat/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatSessionController struct {
	service *services.ChatSessionService
}

func NewChatSessionController(service *services.ChatSessionService) *ChatSessionController {
	return &ChatSessionController{service: service}
}

func (c *ChatSessionController) InitChatSession(ctx *gin.Context) {
	var req dto.CreateChatSessionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if req.Session == "" || req.ProductId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "session and product_id required",
		})
		return
	}

	existing, err := c.service.GetChatSessionByUser(req.Session, req.ProductId)

	if err != nil && err != repository.ErrChatSessionNotFound {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if existing != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "existing_session",
			"token":  existing.Token,
		})
		return
	}

	if err := c.service.InitChatSession(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "session_created",
		"token":  req.Token,
	})

}
