// package consumer

// import (
// 	"encoding/json"
// 	"log"

// 	stockdb "github.com/example/stock-service/internal/db"
// 	pb "github.com/example/stock-service/proto"
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// type CancelledOrder struct {
// 	OrderID string `json:"order_id"`
// 	Items   []struct {
// 		SKU      string `json:"sku"`
// 		Quantity int    `json:"quantity"`
// 	} `json:"items"`
// }

// func StartOrderCancelledConsumer(stockServer pb.StockServiceServer) {
// 	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672/")
// 	if err != nil {
// 		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatalf("❌ Failed to open RabbitMQ channel: %v", err)
// 	}

// 	// Ensure queue exists
// 	_, err = ch.QueueDeclare(
// 		"order_cancelled",
// 		true,  // durable
// 		false, // auto-delete
// 		false, // exclusive
// 		false, // no-wait
// 		nil,
// 	)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to declare order.cancelled queue: %v", err)
// 	}

// 	msgs, err := ch.Consume(
// 		"order_cancelled",
// 		"",
// 		true,  // auto-ack
// 		false, // not exclusive
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to register consumer: %v", err)
// 	}

// 	log.Println("📥 Listening on queue: order_cancelled ...")

// 	// Dependencies
// 	db := stockdb.NewMySQLConnection()

// 	// Start listening
// 	go func() {
// 		for msg := range msgs {
// 			var event CancelledOrder

// 			if err := json.Unmarshal(msg.Body, &event); err != nil {
// 				log.Println("❌ Failed to parse cancelled order event:", err)
// 				continue
// 			}

// 			log.Printf("🚫 Order Cancelled Received: %s\n", event.OrderID)

// 			// for _, item := range event.Items {
// 			// 	// Increase stock for each cancelled item
// 			// 	resp, err := stockServer.ReleaseStock(&pb.ReleaseStockRequest{}item.SKU, item.Quantity)
// 			// 	if err != nil {
// 			// 		log.Printf("❌ Failed to release stock for SKU=%s: %v\n", item.SKU, err)
// 			// 	} else {
// 			// 		log.Printf("↩️ Stock restored for SKU=%s, qty=%d\n", item.SKU, item.Quantity)
// 			// 	}
// 			// }

// 			log.Printf("✔️ Stock restored for cancelled order: %s\n", event.OrderID)
// 		}
// 	}()
// }
