package server

import (
	"context"
	"log"

	"github.com/example/stock-service/internal/service"
	pb "github.com/example/stock-service/proto"
)

type stockServer struct {
	pb.UnimplementedStockServiceServer
	service *service.StockService
}

func NewStockServer(s *service.StockService) *stockServer {
	return &stockServer{
		service: s,
	}
}

// func (s *stockServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
// 	// orderID := req.GetOrderId()
// 	// items := req.GetItems()
// 	fmt.Println("Stock Reserved Successfully")
// 	return &pb.ReserveStockResponse{Success: true}, nil
// }

// func (s *stockServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
// 	// orderID := req.GetOrderId()
// 	// items := req.GetItems()

// 	orderID := req.GetOrderId()
// 	items := req.GetItems()

// 	fmt.Println(items)

// 	// ---- BEGIN TRANSACTION ----
// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return &pb.ReserveStockResponse{
// 			Success: false,
// 			Message: "Failed to start transaction",
// 		}, err
// 	}

// 	for _, item := range items {
// 		sku := item.Sku
// 		qty := int(item.Quantity)

// 		// ---- Check stock row ----
// 		var available, reserved int
// 		err := tx.QueryRow(`SELECT qty_available, qty_reserved FROM stocks WHERE sku=? FOR UPDATE`, sku).
// 			Scan(&available, &reserved)
// 		if err == sql.ErrNoRows {
// 			tx.Rollback()
// 			return &pb.ReserveStockResponse{
// 				Success: false,
// 				Message: fmt.Sprintf("SKU %s not found", sku),
// 			}, nil
// 		}
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}

// 		// ---- Validate availability ----
// 		if available < qty {
// 			tx.Rollback()
// 			return &pb.ReserveStockResponse{
// 				Success:    false,
// 				Message:    fmt.Sprintf("Not enough stock for %s", sku),
// 				FailedSkus: []string{sku},
// 			}, nil
// 		}

// 		// ---- Deduct qty_available and increase reserved ----
// 		_, err = tx.Exec(
// 			`UPDATE stocks
//                SET qty_available = qty_available - ?,
//                    qty_reserved = qty_reserved + ?
//              WHERE sku=?`,
// 			qty, qty, sku,
// 		)
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}

// 		// ---- Insert reservation record ----
// 		_, err = tx.Exec(
// 			`INSERT INTO reservations (order_id, sku, quantity) VALUES (?, ?, ?)`,
// 			orderID, sku, qty,
// 		)
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 	}

// 	// ---- COMMIT TRANSACTION ----
// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}

// 	log.Printf("Reserved stock for Order %s successfully\n", orderID)

// 	return &stockpb.ReserveStockResponse{
// 		Success: true,
// 		Message: "Stock reserved successfully",
// 	}, nil
// }

func (s *stockServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {

	log.Println("ReserveStock called")

	items := make(map[string]int)
	for _, it := range req.Items {
		items[it.Sku] = int(it.Quantity)
	}

	err := s.service.ReserveStock(req.OrderId, items)
	if err != nil {
		return &pb.ReserveStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ReserveStockResponse{
		Success: true,
		Message: "Stock reserved successfully",
	}, nil
}

// func (s *stockServer) ReleaseStock(ctx context.Context, req *stockpb.ReleaseStockRequest) (*stockpb.ReleaseStockResponse, error) {
// 	orderID := req.GetOrderId()

// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Step 1: Get all active reservations for this order
// 	rows, err := tx.Query(`
//         SELECT sku, quantity
//         FROM reservations
//         WHERE order_id = ? AND released = FALSE
//         FOR UPDATE
//     `, orderID)
// 	if err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	type item struct {
// 		sku string
// 		qty int
// 	}
// 	var reservedItems []item

// 	for rows.Next() {
// 		var i item
// 		if err := rows.Scan(&i.sku, &i.qty); err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 		reservedItems = append(reservedItems, i)
// 	}

// 	if len(reservedItems) == 0 {
// 		tx.Rollback()
// 		return &stockpb.ReleaseStockResponse{
// 			Success: false,
// 			Message: "No reservations found",
// 		}, nil
// 	}

// 	// Step 2: Release stock for each reserved SKU
// 	for _, it := range reservedItems {
// 		_, err = tx.Exec(`
//             UPDATE stocks
//             SET qty_available = qty_available + ?,
//                 qty_reserved = qty_reserved - ?
//             WHERE sku = ?
//         `, it.qty, it.qty, it.sku)
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 	}

// 	// Step 3: Mark reservations as released
// 	_, err = tx.Exec(`
//         UPDATE reservations
//         SET released = TRUE
//         WHERE order_id = ?
//     `, orderID)
// 	if err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	// Commit final
// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}

// 	return &stockpb.ReleaseStockResponse{
// 		Success: true,
// 		Message: "Stock released successfully",
// 	}, nil
// }

func (s *stockServer) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReleaseStockResponse, error) {
	log.Println("ReleaseStock called")

	// Convert proto items into map for service layer
	itemMap := make(map[string]int)
	for _, it := range req.Items {
		itemMap[it.Sku] = int(it.Quantity)
	}

	err := s.service.ReleaseStock(req.OrderId, itemMap)
	if err != nil {
		return &pb.ReleaseStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ReleaseStockResponse{
		Success: true,
		Message: "Stock released successfully",
	}, nil
}

// func (s *stockServer) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
// 	sku := req.GetSku()

// 	var available, reserved int32

// 	s.db.QueryRow(
// 		`SELECT qty_available, qty_reserved FROM stocks WHERE sku = ?`, sku).Scan(&available, &reserved)

// 	return &pb.GetStockResponse{Sku: req.Sku, QtyAvailable: available, QtyReserved: reserved}, nil
// }

func (s *stockServer) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {

	available, reserved, err := s.service.GetStock(req.Sku)
	if err != nil {
		return &pb.GetStockResponse{
			Sku: "",
		}, nil
	}

	return &pb.GetStockResponse{
		Sku:          req.Sku,
		QtyAvailable: int32(available),
		QtyReserved:  int32(reserved),
	}, nil
}
