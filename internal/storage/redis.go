package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"test-server-go/internal/config"
)

type Redis struct {
	*redis.Client
}

func NewRedis(ctx context.Context, cfg config.Config) (*Redis, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Ip, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{redisClient}, nil
}
