package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return conn, ch
}
