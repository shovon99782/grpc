package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishOrderCreated(ch *amqp.Channel, event interface{}) error {
	body, _ := json.Marshal(event)

	return ch.PublishWithContext(
		context.Background(),
		"order.events",  // exchange
		"order.created", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
