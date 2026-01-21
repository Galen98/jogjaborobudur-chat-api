package services

import (
	"errors"
	"os"

	pushnotifications "github.com/pusher/push-notifications-go"
)

const AdminChatInterest = "admin-chat"

type AdminPushService struct {
	client pushnotifications.PushNotifications
}

func NewAdminPushService() (*AdminPushService, error) {
	instanceID := os.Getenv("PUSHER_BEAMS_INSTANCE_ID")
	secretKey := os.Getenv("PUSHER_BEAMS_SECRET_KEY")

	if instanceID == "" || secretKey == "" {
		return nil, errors.New("pusher beams config missing")
	}

	client, err := pushnotifications.New(instanceID, secretKey)
	if err != nil {
		return nil, err
	}

	return &AdminPushService{
		client: client,
	}, nil
}

func (s *AdminPushService) NotifyNewChat(
	productName string,
	message string,
) error {

	body := truncate(message, 80)

	_, err := s.client.PublishToInterests(
		[]string{AdminChatInterest},
		map[string]interface{}{
			"web": map[string]interface{}{
				"notification": map[string]interface{}{
					"title": "New Chat: " + productName,
					"body":  body,
				},
				"data": map[string]interface{}{
					"type": "new_chat",
				},
			},
		},
	)

	return err
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
