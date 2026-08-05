package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/utils"
)

var redisClient *redis.Client

func InitRedis(ctx context.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     utils.Get(constant.EnvKeyRedisAddr, "127.0.0.1:6379"),
		Password: utils.Get(constant.EnvKeyRedisPassword),
		DB:       utils.GetAsInt(constant.EnvKeyRedisDB, 0),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return fmt.Errorf("connect redis: %w", err)
	}

	redisClient = client
	return nil
}

func getRedis() (*redis.Client, error) {
	if redisClient == nil {
		return nil, errors.New("redis is not initialized")
	}

	return redisClient, nil
}
