package main

import (
	"log"
	"net/http"

	rabbitmq "github.com/example/analytics-service/internal/consumer"
	elasticsearch "github.com/example/analytics-service/internal/elastic"
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

	// ---- 3. Connect Elasticsearch ONCE ----
	// es, err := elasticsearch.NewClient(elasticsearch.Config{
	// 	Addresses: []string{"http://localhost:9200"},
	// })
	// if err != nil {
	// 	log.Fatalf("❌ Elasticsearch connection failed: %v", err)
	// }

	elasticsearch.InitElasticsearch()

	go rabbitmq.StartConsumer(ch, elasticsearch.ES)

	log.Println("Analytics Service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
