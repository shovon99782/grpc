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
	"google.golang.org/protobuf/types/known/timestamppb"
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

	stockServiceConn, err := grpc.Dial("stock-service:50052", grpc.WithInsecure())
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
		"OrderID":    newOrderId,
		"Status":     "CREATED",
		"CustomerID": req.GetCustomerId(),
		"Items":      req.GetItems(),
		"Time":       time.Now().String(),
	}

	body, _ := json.Marshal(event)

	go func() {
		err = s.rabbit.Publish("order_created", body)
		if err != nil {
			fmt.Println(err)
		}

		log.Println("🎉 Order created and event published")
	}()

	return &orderpb.CreateOrderResponse{OrderId: newOrderId, Status: "CREATED"}, nil
}

func (s *orderServer) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.Empty, error) {
	log.Println("UpdateOrderStatus called")

	query := `UPDATE orders SET status = ?, updated_at = NOW() WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, req.Status, req.OrderId)
	if err != nil {
		log.Printf("❌ Failed to update order: %v", err)
		return nil, err
	}

	if req.Status == "CANCELLED" {
		rows, err := s.db.Query(`
            SELECT sku, quantity FROM order_items WHERE order_id = ?
        `, req.OrderId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var items []*stockpb.ReserveItem
		for rows.Next() {
			var sku string
			var qty int
			rows.Scan(&sku, &qty)

			items = append(items, &stockpb.ReserveItem{
				Sku:      sku,
				Quantity: int32(qty),
			})
		}

		// Build event payload
		event := stockpb.ReleaseStockRequest{
			OrderId: req.OrderId,
			Items:   items,
		}

		body, _ := json.Marshal(event)

		// Publish to RabbitMQ
		err = s.rabbit.Publish("order_cancelled", body)
		if err != nil {
			log.Printf("❌ Failed to publish cancel event: %v", err)
		} else {
			log.Printf("📢 Published order_cancelled event for OrderID=%s", req.OrderId)
		}
	}

	log.Printf("✅ Order %s updated to status %s", req.OrderId, req.Status)

	return &orderpb.Empty{}, nil
}

func (s *orderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	log.Println("GetOrder called")

	// Fetch main order info
	query := `
        SELECT id, customer_id, status, total_amount 
        FROM orders 
        WHERE id = ?
    `
	var (
		orderID, customerID, status string
		total                       float64
	)

	err := s.db.QueryRowContext(ctx, query, req.OrderId).
		Scan(&orderID, &customerID, &status, &total)

	if err != nil {
		log.Printf("❌ Order not found: %v", err)
		return nil, err
	}

	// Fetch order items
	itemsQuery := `
        SELECT sku, quantity, unit_price 
        FROM order_items
        WHERE order_id = ?
    `

	rows, err := s.db.QueryContext(ctx, itemsQuery, req.OrderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*orderpb.OrderItem
	for rows.Next() {
		var sku string
		var qty int32
		var price float64

		rows.Scan(&sku, &qty, &price)

		items = append(items, &orderpb.OrderItem{
			Sku:       sku,
			Quantity:  qty,
			UnitPrice: float64(price),
		})
	}

	return &orderpb.GetOrderResponse{
		OrderId:     orderID,
		CustomerId:  customerID,
		Status:      status,
		TotalAmount: float64(total),
		Items:       items,
	}, nil
}

func (s *orderServer) GetOrdersByCustomer(ctx context.Context, req *orderpb.GetOrdersByCustomerRequest) (*orderpb.GetOrdersByCustomerResponse, error) {
	log.Println("GetOrdersByCustomer called for customer:", req.CustomerId)

	// Query orders for the customer
	const ordersQ = `
		SELECT id, customer_id, total_amount, status, created_at
		FROM orders
		WHERE customer_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, ordersQ, req.CustomerId)
	if err != nil {
		log.Printf("failed to query orders: %v", err)
		return nil, err
	}
	defer rows.Close()

	var resOrders []*orderpb.GetOrderResponse

	for rows.Next() {
		var (
			orderID    string
			customerID string
			status     string
			totalAmt   float64
			createdAt  time.Time
		)

		if err := rows.Scan(&orderID, &customerID, &totalAmt, &status, &createdAt); err != nil {
			log.Printf("failed to scan order row: %v", err)
			return nil, err
		}

		// Query items for this order
		const itemsQ = `
			SELECT sku, quantity, unit_price
			FROM order_items
			WHERE order_id = ?
		`
		itemRows, err := s.db.QueryContext(ctx, itemsQ, orderID)
		if err != nil {
			log.Printf("failed to query order items for order %s: %v", orderID, err)
			return nil, err
		}

		var items []*orderpb.OrderItem
		for itemRows.Next() {
			var (
				sku       string
				qty       int
				unitPrice float64
			)
			if err := itemRows.Scan(&sku, &qty, &unitPrice); err != nil {
				itemRows.Close()
				log.Printf("failed to scan item row for order %s: %v", orderID, err)
				return nil, err
			}
			items = append(items, &orderpb.OrderItem{
				Sku:       sku,
				Quantity:  int32(qty),
				UnitPrice: unitPrice,
			})
		}
		itemRows.Close()

		resOrders = append(resOrders, &orderpb.GetOrderResponse{
			OrderId:     orderID,
			CustomerId:  customerID,
			Items:       items,
			TotalAmount: totalAmt,
			Status:      status,
			CreatedAt:   timestamppb.New(createdAt),
		})
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows error: %v", err)
		return nil, err
	}

	return &orderpb.GetOrdersByCustomerResponse{
		Orders: resOrders,
	}, nil
}
