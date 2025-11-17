package consumer

import (
	"encoding/json"
	"log"

	service "github.com/example/stock-service/internal/service"
	pb "github.com/example/stock-service/proto"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderCancelConsumer struct {
	ch      *amqp.Channel
	service *service.StockService
}

func NewOrderCancelConsumer(ch *amqp.Channel, s *service.StockService) *OrderCancelConsumer {
	return &OrderCancelConsumer{
		ch:      ch,
		service: s,
	}
}

func (c *OrderCancelConsumer) Start() error {
	// ensure queue exists
	_, err := c.ch.QueueDeclare(
		"order_cancelled",
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.ch.Consume(
		"order_cancelled",
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("🔁 OrderCancelConsumer: Listening for cancelled orders...")

	go func() {
		for msg := range msgs {
			var req pb.ReleaseStockRequest

			// Parse incoming cancel event
			if err := json.Unmarshal(msg.Body, &req); err != nil {
				log.Printf("❌ Failed to parse order_cancelled message: %v", err)
				continue
			}

			log.Printf("🛑 Order cancelled → releasing stock: %s", req.OrderId)

			// Convert repeated items → map[string]int
			skuQty := make(map[string]int)
			for _, item := range req.Items {
				skuQty[item.Sku] = int(item.Quantity)
			}

			// Call the internal service layer
			err := c.service.ReleaseStock(req.OrderId, skuQty)
			if err != nil {
				log.Printf("❌ Failed to release stock for OrderID=%s: %v", req.OrderId, err)
				continue
			}

			log.Printf("✅ Stock successfully released for OrderID=%s", req.OrderId)
		}
	}()

	return nil
}
