package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/example/order-service/internal/rabbitmq"
	orderpb "github.com/example/order-service/proto/order"
	stockpb "github.com/example/order-service/proto/stock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService struct {
	db     *sql.DB
	rabbit *rabbitmq.RabbitMQ
}

func NewOrderService(db *sql.DB, rabbit *rabbitmq.RabbitMQ) *OrderService {
	return &OrderService{db: db, rabbit: rabbit}
}

func calculateTotal(items []*orderpb.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}

func (s *OrderService) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	orderID := uuid.New().String()

	// --- 1️⃣ Reserve Stock ---
	stockItems := []*stockpb.ReserveItem{}
	for _, item := range req.Items {
		stockItems = append(stockItems, &stockpb.ReserveItem{
			Sku:      item.Sku,
			Quantity: item.Quantity,
		})
	}

	stockConn, err := grpc.Dial("stock-service:50052", grpc.WithInsecure())
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stockConn.Close()

	stockClient := stockpb.NewStockServiceClient(stockConn)

	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = stockClient.ReserveStock(ctx2, &stockpb.ReserveStockRequest{
		OrderId: orderID,
		Items:   stockItems,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// --- 2️⃣ Insert Order ---
	_, err = tx.Exec(
		`INSERT INTO orders (id, customer_id, total_amount, status) VALUES (?, ?, ?, 'CREATED')`,
		orderID, req.CustomerId, calculateTotal(req.Items),
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// --- 3️⃣ Insert Items ---
	for _, item := range req.Items {
		_, err = tx.Exec(
			`INSERT INTO order_items (order_id, sku, quantity, unit_price) VALUES (?, ?, ?, ?)`,
			orderID, item.Sku, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	// --- 4️⃣ Publish Event ---
	event := map[string]interface{}{
		"OrderID": orderID,
		"Status":  "CREATED",
		"Items":   req.Items,
		"Time":    time.Now().String(),
	}

	body, _ := json.Marshal(event)

	go s.rabbit.Publish("order_created", body)

	return &orderpb.CreateOrderResponse{
		OrderId: orderID,
		Status:  "CREATED",
	}, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.Empty, error) {

	_, err := s.db.Exec(
		`UPDATE orders SET status = ?, updated_at = NOW() WHERE id = ?`,
		req.Status, req.OrderId,
	)
	if err != nil {
		return nil, err
	}

	// Handle cancel
	if req.Status == "CANCELLED" {
		rows, err := s.db.Query(`SELECT sku, quantity FROM order_items WHERE order_id = ?`, req.OrderId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var items []*stockpb.ReserveItem
		for rows.Next() {
			var sku string
			var qty int
			rows.Scan(&sku, &qty)
			items = append(items, &stockpb.ReserveItem{Sku: sku, Quantity: int32(qty)})
		}

		event := stockpb.ReleaseStockRequest{
			OrderId: req.OrderId,
			Items:   items,
		}

		body, _ := json.Marshal(event)

		s.rabbit.Publish("order_cancelled", body)
	}

	return &orderpb.Empty{}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {

	var (
		id, customer, status string
		total                float64
	)

	err := s.db.QueryRow(`
		SELECT id, customer_id, status, total_amount 
		FROM orders WHERE id = ?`, req.OrderId).
		Scan(&id, &customer, &status, &total)

	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(`
		SELECT sku, quantity, unit_price 
		FROM order_items WHERE order_id = ?`, req.OrderId)
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
			UnitPrice: price,
		})
	}

	return &orderpb.GetOrderResponse{
		OrderId:     id,
		CustomerId:  customer,
		Status:      status,
		TotalAmount: total,
		Items:       items,
	}, nil
}

func (s *OrderService) GetOrdersByCustomer(ctx context.Context, req *orderpb.GetOrdersByCustomerRequest) (*orderpb.GetOrdersByCustomerResponse, error) {

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

	var resultOrders []*orderpb.GetOrderResponse

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

		// Fetch items for order
		const itemsQ = `
			SELECT sku, quantity, unit_price
			FROM order_items
			WHERE order_id = ?
		`

		itemsRows, err := s.db.QueryContext(ctx, itemsQ, orderID)
		if err != nil {
			log.Printf("failed to query order items: %v", err)
			return nil, err
		}

		var items []*orderpb.OrderItem

		for itemsRows.Next() {
			var (
				sku       string
				qty       int
				unitPrice float64
			)
			if err := itemsRows.Scan(&sku, &qty, &unitPrice); err != nil {
				itemsRows.Close()
				log.Printf("failed to scan item row: %v", err)
				return nil, err
			}

			items = append(items, &orderpb.OrderItem{
				Sku:       sku,
				Quantity:  int32(qty),
				UnitPrice: unitPrice,
			})
		}

		itemsRows.Close()

		resultOrders = append(resultOrders, &orderpb.GetOrderResponse{
			OrderId:     orderID,
			CustomerId:  customerID,
			Items:       items,
			TotalAmount: totalAmt,
			Status:      status,
			CreatedAt:   timestamppb.New(createdAt),
		})
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
		return nil, err
	}

	return &orderpb.GetOrdersByCustomerResponse{Orders: resultOrders}, nil
}
