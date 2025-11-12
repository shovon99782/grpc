package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	orderpb "github.com/example/order-service/proto/order"
)

func main() {
	// Connect to OrderService running at localhost:50051
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer conn.Close()

	client := orderpb.NewOrderServiceClient(conn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare CreateOrder request
	req := &orderpb.CreateOrderRequest{
		CustomerId: "cust-123",
		Items: []*orderpb.OrderItem{
			{Sku: "SKU123", Quantity: 2, UnitPrice: 49.99},
			{Sku: "SKU456", Quantity: 1, UnitPrice: 29.99},
		},
	}

	// Call CreateOrder via gRPC
	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		log.Fatalf("CreateOrder failed: %v", err)
	}

	fmt.Printf("✅ Order Created Successfully!\nOrder ID: %s\nStatus: %s\n", resp.OrderId, resp.Status)
}
