package rabbitmq

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderCreatedEvent struct {
	OrderID string      `json:"order_id"`
	Status  string      `json:"status"`
	Items   interface{} `json:"items"`
	Time    string      `json:"time"`
}

// placeholder consumer logic: connect to RabbitMQ, consume order.events, index into ES
func StartConsumer(ch *amqp.Channel, es *elasticsearch.Client) {
	_, err := ch.QueueDeclare(
		"order_created",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to declare queue: %v", err)
	}

	// Start consuming
	msgs, err := ch.Consume(
		"order_created",
		"analytics-consumer",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to start consumer: %v", err)
	}

	log.Println("📡 Listening for events on queue: order_created")

	// Handle messages
	go func() {
		for m := range msgs {

			var event OrderCreatedEvent
			err := json.Unmarshal(m.Body, &event)
			if err != nil {
				log.Println("❌ Failed to decode event:", err)
				continue
			}

			log.Println("📥 Received OrderCreated Event:")
			log.Println("   OrderID:", event.OrderID)
			log.Println("   Status:", event.Status)
			log.Println("   Time:", event.Time)
			log.Println("   Items:", event.Items)
			log.Println("-------------------------------------------")

			// TODO:
			//   - Send to Elasticsearch
			//   - Enrich with customer data
			//   - Build analytics dashboard
			docBytes, _ := json.Marshal(event)
			res, err := es.Index(
				"orders",
				bytes.NewReader(docBytes),
				es.Index.WithDocumentID(event.OrderID),
				es.Index.WithRefresh("true"),
				es.Index.WithContext(context.Background()),
			)
			if err != nil {
				log.Printf("❌ Elasticsearch insert error: %v", err)
				continue
			}
			res.Body.Close()

			log.Printf("✅ Order %s inserted into Elasticsearch", event.OrderID)
		}
	}()

}
