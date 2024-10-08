package rabbitmq

import (
	"fmt"
	"ticket-app/config"

	"github.com/streadway/amqp"
)

type RabbitMQConnection struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

var cfg = config.LoadRabbitMQConfig()

func ConnectRabbitMQ() (*RabbitMQConnection, error) {
	rabbitMQURL := cfg.RabbitMQURL

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQConnection{Conn: conn, Ch: ch}, nil
}

func (r *RabbitMQConnection) Close() {
	if r.Ch != nil {
		r.Ch.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}
