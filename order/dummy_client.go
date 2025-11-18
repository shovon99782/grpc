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
	// Connect to OrderService
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect to order service: %v", err)
	}
	defer conn.Close()

	client := orderpb.NewOrderServiceClient(conn)

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//-------------------------------------------
	// 1️⃣ CREATE ORDER
	//-------------------------------------------
	req := &orderpb.CreateOrderRequest{
		CustomerId: "cust-456",
		Items: []*orderpb.OrderItem{
			{Sku: "SKU123", Quantity: 2, UnitPrice: 49.99},
			{Sku: "SKU456", Quantity: 1, UnitPrice: 29.99},
		},
	}

	createResp, err := client.CreateOrder(ctx, req)
	if err != nil {
		log.Fatalf("CreateOrder failed: %v", err)
	}

	fmt.Println("=== CreateOrder Response ===")
	fmt.Printf("Order ID: %s\n", createResp.OrderId)
	fmt.Printf("Status: %s\n\n", createResp.Status)

	orderID := createResp.OrderId
	fmt.Println("NEW ORDER ID:", orderID)
	// orderID := "550e889b-1929-4c3e-9714-831ceeaff150"
	//-------------------------------------------
	// 2️⃣ UPDATE ORDER STATUS
	//-------------------------------------------
	updateReq := &orderpb.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  "CANCELLED",
	}

	_, err = client.UpdateOrderStatus(ctx, updateReq)
	if err != nil {
		log.Fatalf("UpdateOrderStatus failed: %v", err)
	}

	fmt.Println("=== UpdateOrderStatus Response ===")
	fmt.Printf("Order %s status updated to CONFIRMED\n\n", orderID)

	// //-------------------------------------------
	// // 3️⃣ GET ORDER BY ID
	// //-------------------------------------------
	getReq := &orderpb.GetOrderRequest{
		OrderId: orderID,
	}

	getResp, err := client.GetOrder(ctx, getReq)
	if err != nil {
		log.Fatalf("GetOrder failed: %v", err)
	}

	fmt.Println("=== GetOrder Response ===")
	fmt.Printf("Order ID: %s\n", getResp.OrderId)
	fmt.Printf("Customer ID: %s\n", getResp.CustomerId)
	fmt.Printf("Status: %s\n", getResp.Status)
	fmt.Printf("Total Amount: %.2f\n", getResp.TotalAmount)
	fmt.Printf("Created At: %v\n", getResp.CreatedAt.AsTime())

	fmt.Println("Items:")
	for _, item := range getResp.Items {
		fmt.Printf("  - SKU: %s Qty: %d Price: %.2f\n",
			item.Sku, item.Quantity, item.UnitPrice)
	}
	fmt.Println()

	//-------------------------------------------
	// 4️⃣ GET ORDERS BY CUSTOMER
	//-------------------------------------------
	custReq := &orderpb.GetOrdersByCustomerRequest{
		CustomerId: "cust-123",
	}

	custResp, err := client.GetOrdersByCustomer(ctx, custReq)
	if err != nil {
		log.Fatalf("GetOrdersByCustomer failed: %v", err)
	}

	fmt.Println("=== GetOrdersByCustomer Response ===")
	for _, order := range custResp.Orders {
		fmt.Printf("OrderID: %s | Amount: %.2f | Status: %s | Created: %v\n",
			order.OrderId,
			order.TotalAmount,
			order.Status,
			order.CreatedAt.AsTime(),
		)
	}
}
