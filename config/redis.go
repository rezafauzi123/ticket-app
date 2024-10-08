package config

import (
	"os"
	"ticket-app/pkg/log"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	RedisHost string
	RedisPort string
}

func LoadRedisConfig() *RedisConfig {
	logger := log.GetLogger()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	return &RedisConfig{
		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
	}
}
