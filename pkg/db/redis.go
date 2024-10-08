package db

import (
	"context"
	"ticket-app/config"
	"ticket-app/pkg/log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var cfg = config.LoadRedisConfig()

func InitRedis() {
	logger := log.GetLogger()
	RedisClient = redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost + ":" + cfg.RedisPort,
		DB:   0,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal(err)
	}
}
