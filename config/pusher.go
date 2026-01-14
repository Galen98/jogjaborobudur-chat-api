package config

import (
	"os"

	"github.com/pusher/pusher-http-go"
)

var Pusher *pusher.Client

func InitPusher() {
	Pusher = &pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
		Secure:  true,
	}
}
