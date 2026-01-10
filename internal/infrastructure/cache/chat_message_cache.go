package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"jogjaborobudur-chat/internal/domain/chat/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type ChatMessageCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewChatMessageCache(client *redis.Client) *ChatMessageCache {
	return &ChatMessageCache{
		client: client,
		ttl:    24 * time.Hour,
	}
}

func (c *ChatMessageCache) key(token string) string {
	return fmt.Sprintf("chat:messages:%s", token)
}

func (c *ChatMessageCache) Get(token string) (*entity.ChatConversation, error) {
	val, err := c.client.Get(ctx, c.key(token)).Result()
	if err != nil {
		return nil, err
	}

	var conv entity.ChatConversation
	if err := json.Unmarshal([]byte(val), &conv); err != nil {
		return nil, err
	}
	return &conv, nil
}

func (c *ChatMessageCache) Set(conv *entity.ChatConversation) error {
	data, err := json.Marshal(conv)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.key(conv.Token), data, c.ttl).Err()
}

func (c *ChatMessageCache) PushMessage(token string, msg entity.ChatData) error {
	key := c.key(token)
	data, _ := json.Marshal(msg)
	pipe := c.client.TxPipeline()

	pipe.RPush(ctx, key, data)

	pipe.LTrim(ctx, key, -20, -1)

	pipe.Expire(ctx, key, c.ttl)

	_, err := pipe.Exec(ctx)
	return err
}
