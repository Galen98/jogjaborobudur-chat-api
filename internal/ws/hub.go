package ws

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMessage
	redisSub   *redis.Client
}

type BroadcastMessage struct {
	Token string
	Data  []byte
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMessage),
		redisSub:   redisClient,
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			if _, ok := h.rooms[client.Token]; !ok {
				h.rooms[client.Token] = make(map[*Client]bool)
			}
			h.rooms[client.Token][client] = true

		case client := <-h.unregister:
			if room, ok := h.rooms[client.Token]; ok {

				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.Send)
				}

				if len(room) == 0 {
					delete(h.rooms, client.Token)
				}
			}

		case msg := <-h.broadcast:
			if room, ok := h.rooms[msg.Token]; ok {

				for c := range room {
					select {
					case c.Send <- msg.Data:
					default:
						close(c.Send)
						delete(room, c)
					}
				}
			}
		}
	}
}

func (h *Hub) ListenRedis() {
	ctx := context.Background()

	sub := h.redisSub.PSubscribe(
		ctx,
		"chat:*",
		"session:*",
		"admin:sessions",
	)

	for msg := range sub.Channel() {
		h.broadcast <- BroadcastMessage{
			Token: msg.Channel,
			Data:  []byte(msg.Payload),
		}
	}
}

func (h *Hub) Broadcast(msg BroadcastMessage) {
	h.broadcast <- msg
}

func (h *Hub) Register(c *Client) {
	h.register <- c
}

func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}
