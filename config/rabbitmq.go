package config

import (
	"os"
	"ticket-app/pkg/log"

	"github.com/joho/godotenv"
)

type RabbitMQConfig struct {
	RabbitMQURL        string
	RabbitMQExchange   string
	RabbitMQQueue      string
	RabbitMQRoutingKey string
}

func LoadRabbitMQConfig() *RabbitMQConfig {
	logger := log.GetLogger()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	return &RabbitMQConfig{
		RabbitMQURL:        os.Getenv("RABBITMQ_URL"),
		RabbitMQExchange:   os.Getenv("RABBITMQ_EXCHANGE"),
		RabbitMQQueue:      os.Getenv("RABBITMQ_QUEUE"),
		RabbitMQRoutingKey: os.Getenv("RABBITMQ_ROUTING_KEY"),
	}
}
