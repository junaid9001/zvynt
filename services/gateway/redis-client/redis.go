package redisclient

import (
	"context"
	"fmt"
	"log"

	"github.com/junaid9001/zvynt/gateway/config"
	"github.com/redis/go-redis/v9"
)

func RedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.REDIS_ADDR,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	fmt.Println("redis connected:", pong)

	return rdb
}
