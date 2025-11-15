package main

import (
	"log"
	"net/http"

	rabbitmq "github.com/example/analytics-service/internal/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// start a minimal HTTP server for search API (stub)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to create channel: %v", err)
	}
	log.Println("✅ RabbitMQ channel opened")
	defer ch.Close()

	go rabbitmq.StartConsumer(ch)

	log.Println("Analytics Service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
