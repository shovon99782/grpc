package rabbitmq

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn *amqp.Connection
}

func NewRabbitMQ(url string) *RabbitMQ {
	// conn, err := amqp.Dial(url)
	// if err != nil {
	// 	log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	// }

	var conn *amqp.Connection
	var err error
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}

		log.Printf("Retrying RabbitMQ in 3s... (%v)", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ after retries: %v", err)
	}

	log.Println("🐰 RabbitMQ connected successfully")
	return &RabbitMQ{Conn: conn}
}

func (r *RabbitMQ) DeclareQueue(queue string) (*amqp.Channel, error) {
	ch, err := r.Conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	return ch, nil
}
