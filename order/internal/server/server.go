package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/example/order-service/internal/rabbitmq"
	orderpb "github.com/example/order-service/proto/order"
	stockpb "github.com/example/order-service/proto/stock"
	"google.golang.org/grpc"
)

type orderServer struct {
	orderpb.UnimplementedOrderServiceServer
	rabbit *rabbitmq.RabbitMQ
}

func NewOrderServer(rabbit *rabbitmq.RabbitMQ) *orderServer {

	return &orderServer{
		rabbit: rabbit,
	}
}

func (s *orderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Println("CreateOrder called - stub")
	// TODO: validate, call StockService, save to DB, publish event
	items := req.GetItems()
	for _, item := range items {
		fmt.Println("/////////")
		fmt.Println(item.Sku)
	}

	stockServiceConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer stockServiceConn.Close()

	client := stockpb.NewStockServiceClient(stockServiceConn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stockReq := &stockpb.ReserveStockRequest{
		OrderId: "O-123",
		Items:   []*stockpb.ReserveItem{},
	}

	resp, err := client.ReserveStock(ctx, stockReq)
	if err != nil {
		log.Fatalf("ReserveStock failed: %v", err)
	}

	fmt.Printf("✅ Stock Reserved Successfully!\nSuccess Status: %s\nMessage: %s\n", resp.Success, resp.Message)

	event := map[string]interface{}{
		"order_id": "O-123",
		"status":   "CREATED",
		"items":    req.GetItems(),
		"time":     time.Now().String(),
	}

	body, _ := json.Marshal(event)

	err = s.rabbit.Publish("order_created", body)
	if err != nil {
		return nil, err
	}

	log.Println("🎉 Order created and event published")

	return &orderpb.CreateOrderResponse{OrderId: "stub-order-id", Status: "CREATED"}, nil
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
