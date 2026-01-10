package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedis() *redis.Client {
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))

	opt := &redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + strconv.Itoa(port),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	if os.Getenv("REDIS_SSL") == "True" {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	rdb := redis.NewClient(opt)

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("redis not connect " + err.Error())
	}

	fmt.Println("conect success")

	return rdb
}
