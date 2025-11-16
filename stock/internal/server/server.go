package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "github.com/example/stock-service/proto"
	stockpb "github.com/example/stock-service/proto"
)

type stockServer struct {
	db *sql.DB
	pb.UnimplementedStockServiceServer
}

func NewStockServer(db *sql.DB) *stockServer {
	return &stockServer{
		db: db,
	}
}

// func (s *stockServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
// 	// orderID := req.GetOrderId()
// 	// items := req.GetItems()
// 	fmt.Println("Stock Reserved Successfully")
// 	return &pb.ReserveStockResponse{Success: true}, nil
// }

func (s *stockServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	// orderID := req.GetOrderId()
	// items := req.GetItems()

	orderID := req.GetOrderId()
	items := req.GetItems()

	fmt.Println(items)

	// ---- BEGIN TRANSACTION ----
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return &pb.ReserveStockResponse{
			Success: false,
			Message: "Failed to start transaction",
		}, err
	}

	for _, item := range items {
		sku := item.Sku
		qty := int(item.Quantity)

		// ---- Check stock row ----
		var available, reserved int
		err := tx.QueryRow(`SELECT qty_available, qty_reserved FROM stocks WHERE sku=? FOR UPDATE`, sku).
			Scan(&available, &reserved)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return &pb.ReserveStockResponse{
				Success: false,
				Message: fmt.Sprintf("SKU %s not found", sku),
			}, nil
		}
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// ---- Validate availability ----
		if available < qty {
			tx.Rollback()
			return &pb.ReserveStockResponse{
				Success:    false,
				Message:    fmt.Sprintf("Not enough stock for %s", sku),
				FailedSkus: []string{sku},
			}, nil
		}

		// ---- Deduct qty_available and increase reserved ----
		_, err = tx.Exec(
			`UPDATE stocks 
               SET qty_available = qty_available - ?, 
                   qty_reserved = qty_reserved + ?
             WHERE sku=?`,
			qty, qty, sku,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// ---- Insert reservation record ----
		_, err = tx.Exec(
			`INSERT INTO reservations (order_id, sku, quantity) VALUES (?, ?, ?)`,
			orderID, sku, qty,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// ---- COMMIT TRANSACTION ----
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	log.Printf("Reserved stock for Order %s successfully\n", orderID)

	return &stockpb.ReserveStockResponse{
		Success: true,
		Message: "Stock reserved successfully",
	}, nil
}

func (s *stockServer) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReserveStockResponse, error) {
	log.Println("ReleaseStock called - stub")
	return &pb.ReserveStockResponse{Success: true}, nil
}

func (s *stockServer) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	log.Println("GetStock called - stub")
	return &pb.GetStockResponse{Sku: req.Sku, QtyAvailable: 100, QtyReserved: 0}, nil
}
