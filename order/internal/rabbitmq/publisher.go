package rabbitmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Publish(queue string, body []byte) error {
	ch, err := r.DeclareQueue(queue) // same connection, new channel
	if err != nil {
		return err
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Failed to publish message:", err)
		return err
	}

	log.Println("Published to queue:", queue)
	return nil
}
