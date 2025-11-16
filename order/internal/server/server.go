package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/example/order-service/internal/rabbitmq"
	orderpb "github.com/example/order-service/proto/order"
	stockpb "github.com/example/order-service/proto/stock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func calculateTotal(items []*orderpb.OrderItem) float64 {
	var total float64 = 0

	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}

	return total
}

type orderServer struct {
	orderpb.UnimplementedOrderServiceServer
	rabbit *rabbitmq.RabbitMQ
	db     *sql.DB
}

func NewOrderServer(db *sql.DB, rabbit *rabbitmq.RabbitMQ) *orderServer {

	return &orderServer{
		rabbit: rabbit,
		db:     db,
	}
}

func (s *orderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Println("CreateOrder called - stub")
	// TODO: validate, call StockService, save to DB, publish event

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	items := req.GetItems()
	stockItems := []*stockpb.ReserveItem{}

	// for _, item := range req.GetItems() {
	// 	stockItems = append(stockItems, &stockpb.ReserveItem{
	// 		Sku:      item.Sku,
	// 		Quantity: item.Quantity,
	// 	})
	// }
	for _, item := range items {
		fmt.Println("/////////")
		fmt.Println(item.Sku)
		stockItems = append(stockItems, &stockpb.ReserveItem{
			Sku:      item.Sku,
			Quantity: item.Quantity,
		})
	}

	stockServiceConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer stockServiceConn.Close()

	client := stockpb.NewStockServiceClient(stockServiceConn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newOrderId := uuid.New().String()

	stockReq := &stockpb.ReserveStockRequest{
		OrderId: newOrderId,
		Items:   stockItems,
	}

	resp, err := client.ReserveStock(ctx, stockReq)
	if err != nil {
		log.Fatalf("ReserveStock failed: %v", err)
	}

	fmt.Printf("✅ Stock Reserved Successfully!\nSuccess Status: %s\nMessage: %s\n", resp.Success, resp.Message)

	// Insert order
	_, err = tx.Exec(
		`INSERT INTO orders (id, customer_id, total_amount, status) VALUES (?, ?, ?, 'CREATED')`,
		newOrderId,
		req.CustomerId,
		calculateTotal(req.Items),
	)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return nil, err
	}

	// Insert items
	for _, item := range req.Items {
		_, err = tx.Exec(
			`INSERT INTO order_items (order_id, sku, quantity, unit_price) VALUES (?, ?, ?, ?)`,
			newOrderId,
			item.Sku,
			item.Quantity,
			item.UnitPrice,
		)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	event := map[string]interface{}{
		"OrderID": newOrderId,
		"Status":  "CREATED",
		"Items":   req.GetItems(),
		"Time":    time.Now().String(),
	}

	body, _ := json.Marshal(event)

	err = s.rabbit.Publish("order_created", body)
	if err != nil {
		return nil, err
	}

	log.Println("🎉 Order created and event published")

	return &orderpb.CreateOrderResponse{OrderId: newOrderId, Status: "CREATED"}, nil
}

func (s *orderServer) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.Empty, error) {
	log.Println("UpdateOrderStatus called - stub")
	return &orderpb.Empty{}, nil
}

func (s *orderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	log.Println("GetOrder called - stub")
	return &orderpb.GetOrderResponse{OrderId: req.OrderId, CustomerId: "cust-stub", Status: "CREATED"}, nil
}

func (s *orderServer) GetOrdersByCustomer(ctx context.Context, req *orderpb.GetOrdersByCustomerRequest) (*orderpb.GetOrdersByCustomerResponse, error) {
	log.Println("GetOrdersByCustomer called - stub")
	return &orderpb.GetOrdersByCustomerResponse{}, nil
}
