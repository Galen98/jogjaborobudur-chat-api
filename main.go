package main

import (
	"jogjaborobudur-chat/config"
	"log"
	"time"

	"jogjaborobudur-chat/internal/domain/chat/email"
	"jogjaborobudur-chat/internal/domain/chat/repository"
	"jogjaborobudur-chat/internal/domain/chat/services"
	httpRouter "jogjaborobudur-chat/internal/http"
	"jogjaborobudur-chat/internal/infrastructure/cache"
	"jogjaborobudur-chat/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "hello"})
	})
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	//pusher
	config.InitPusher()

	// test redis
	rdb := config.NewRedis()
	_ = rdb

	//init db
	if err := config.InitDB(); err != nil {
		log.Fatal("failed connect db", err)
	}

	db := config.DB
	chatSessionRepo := repository.NewChatSessionRepository(db)
	chatMessageRepo := repository.NewChatDataRepository(db)
	// wsHub := ws.NewHub()
	// go wsHub.Run()

	chatCache := cache.NewChatMessageCache(rdb)
	chatDataService := services.NewChatDataService(
		chatSessionRepo,
		chatMessageRepo,
		chatCache,
	)
	userChatRepo := repository.NewUserChatRepository(db)
	chatSessionService := services.NewChatSessionService(chatSessionRepo)
	userChatService := services.NewChatService(userChatRepo)

	smtpCfg := config.LoadSMTPConfig()

	smtpClient := email.NewSMTPClient(
		smtpCfg.Host,
		smtpCfg.Port,
		smtpCfg.Username,
		smtpCfg.Password,
		smtpCfg.From,
	)

	emailService := email.NewEmailService(smtpClient)
	uc := usecase.NewChatUseCase(
		chatDataService,
		chatSessionService,
		userChatService,
		emailService,
	)

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8000",
			"http://localhost:8080",
			"http://localhost:5173",
			"https://jogjaborobudur.com",
			"https://www.jogjaborobudur.com",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ===== Router =====
	httpRouter.SetupRoute(r, db, uc)

	log.Println("server runnig")
	r.Run(":8080")
}
