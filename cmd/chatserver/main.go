package main

import (
	"jogjaborobudur-chat/config"
	"log"

	"jogjaborobudur-chat/internal/domain/chat/email"
	"jogjaborobudur-chat/internal/domain/chat/repository"
	"jogjaborobudur-chat/internal/domain/chat/services"
	httpRouter "jogjaborobudur-chat/internal/http"
	"jogjaborobudur-chat/internal/infrastructure/cache"
	"jogjaborobudur-chat/internal/usecase"
	"jogjaborobudur-chat/internal/ws"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

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
	wsHub := ws.NewHub()
	go wsHub.Run()

	chatCache := cache.NewChatMessageCache(rdb)
	chatDataService := services.NewChatDataService(
		chatSessionRepo,
		chatMessageRepo,
		wsHub,
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
		wsHub,
	)

	// ===== Router =====
	r := httpRouter.SetupRoute(db, uc)

	log.Println("server runnig")
	r.Run(":8080")
}
