package controller

import (
	"jogjaborobudur-chat/internal/domain/chat/dto"
	"jogjaborobudur-chat/internal/domain/chat/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserChatController struct {
	service *services.UserChatService
}

func NewUserChatController(service *services.UserChatService) *UserChatController {
	return &UserChatController{service: service}
}

func (c *UserChatController) CreateUserChat(ctx *gin.Context) {
	var req dto.CreateUserChatRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	resp, err := c.service.CreateUser(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user chat created",
		"data":    resp,
	})
}

func (c *UserChatController) GetUserBySession(ctx *gin.Context) {
	session := ctx.Query("session")
	if session == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "session is required",
		})
		return
	}
	resp, err := c.service.GetBySession(session)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

func (c *UserChatController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "email is required",
		})
		return
	}
	resp, err := c.service.GetByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

func (c *UserChatController) CheckExpiredSession(ctx *gin.Context) {
	var req dto.UserChatRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid request",
		})
		return
	}

	if req.Session == nil {
		ctx.JSON(400, gin.H{
			"message": "session required",
		})
		return
	}

	expired, err := c.service.CheckExpired(*req.Session)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"expired": expired,
	})
}

func (c *UserChatController) DeleteExpiredUsersSheduler(ctx *gin.Context) {
	if err := c.service.DeleteExpiredUsers(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "expired users deleted successfully",
	})
}
