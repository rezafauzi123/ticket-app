package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func PublishMessage(r *RabbitMQConnection, message []byte, routingKey string) error {
	exchange := cfg.RabbitMQExchange
	// routingKey := cfg.RabbitMQRoutingKey

	// Declare an exchange if it does not exist
	err := r.Ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Publish a message
	err = r.Ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Message published: %s", message)
	return nil
}
