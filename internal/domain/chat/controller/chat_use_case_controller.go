package controller

import (
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"jogjaborobudur-chat/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatUseCaseController struct {
	usecase *usecase.ChatUseCase
}

func NewChatUseCaseController(u *usecase.ChatUseCase) *ChatUseCaseController {
	return &ChatUseCaseController{usecase: u}
}

// func (c *ChatUseCaseController) ConnectWS(ctx *gin.Context) {
// 	token := ctx.Param("token")
// 	types := ctx.Param("types")

// 	_ = types
// 	q := ctx.Request.URL.Query()
// 	q.Set("token", token)
// 	ctx.Request.URL.RawQuery = q.Encode()

// 	c.usecase.Hub.ServeWS(ctx.Writer, ctx.Request)
// }

func (c *ChatUseCaseController) SendMessage(ctx *gin.Context) {
	var req dto.SendChatRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	msg, err := c.usecase.SendMessage(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "message_sent",
		"message": msg,
	})

}

func MapChatDataToResponse(cd entity.ChatData) dto.ChatDataResponse {
	return dto.ChatDataResponse{
		ID:               cd.ID,
		ChatSessionToken: cd.ChatSessionToken,
		Message:          cd.Message,
		SenderType:       cd.SenderType,
		Time:             cd.Time.UTC().String(),
		CreatedAt:        cd.CreatedAt.UTC().String(),
		UpdatedAt:        cd.UpdatedAt.UTC().String(),
	}
}

func (c *ChatUseCaseController) GetMessages(ctx *gin.Context) {
	token := ctx.Query("token")
	limitStr := ctx.Query("limit")
	offsetStr := ctx.Query("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if limit == 0 {
		limit = 20
	}

	messages, err := c.usecase.GetMessagesByToken(token, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	responses := make([]dto.ChatDataResponse, 0)
	for _, msg := range messages {
		responses = append(responses, MapChatDataToResponse(msg))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": responses,
	})

}

func (c *ChatUseCaseController) OpenChatByUser(ctx *gin.Context) {
	token := ctx.Param("token")
	types := ctx.Param("types")
	if err := c.usecase.OpenChatByUser(token, types); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "chat_marked_as_read_by_user",
	})
}

func (c *ChatUseCaseController) GetUserSession(ctx *gin.Context) {
	sessionId := ctx.Query("session")
	productIdStr := ctx.Query("product_id")

	productId, _ := strconv.Atoi(productIdStr)

	session, err := c.usecase.GetUserSession(sessionId, uint(productId))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"session": session,
	})
}

func (c *ChatUseCaseController) GetAdminSessions(ctx *gin.Context) {
	sessions, err := c.usecase.GetAdminSessions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

func (c *ChatUseCaseController) GetAllUserChatSession(ctx *gin.Context) {
	sessionId := ctx.Query("session")

	chatsession, err := c.usecase.GetUserChatSessionByUser(sessionId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"session": chatsession,
	})
}
