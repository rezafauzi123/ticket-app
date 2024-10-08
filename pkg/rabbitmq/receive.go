package rabbitmq

import (
	"context"
	"fmt"
	"log"
)

func ReceiveMessages(ctx context.Context, r *RabbitMQConnection, routingKey string, handler func(context.Context, []byte)) error {
	queueName := cfg.RabbitMQQueue
	exchange := cfg.RabbitMQExchange
	// routingKey := cfg.RabbitMQRoutingKey

	// Declare the queue
	q, err := r.Ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Bind the queue to the exchange with a routing key
	err = r.Ch.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %w", err)
	}

	// Start receiving messages
	msgs, err := r.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Handle incoming messages in a goroutine
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// Process the message here, e.g., update ticket status
			handler(ctx, d.Body)
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")
	select {} // Block the function to keep listening for messages
}
