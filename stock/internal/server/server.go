package server

import (
	"context"
	"log"

	pb "github.com/example/stock-service/proto"
)

type stockServer struct {
	pb.UnimplementedStockServiceServer
}

func NewStockServer() *stockServer {
	return &stockServer{}
}

func (s *stockServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	// orderID := req.GetOrderId()
	// items := req.GetItems()

	return &pb.ReserveStockResponse{Success: true}, nil
}

func (s *stockServer) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReserveStockResponse, error) {
	log.Println("ReleaseStock called - stub")
	return &pb.ReserveStockResponse{Success: true}, nil
}

func (s *stockServer) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	log.Println("GetStock called - stub")
	return &pb.GetStockResponse{Sku: req.Sku, QtyAvailable: 100, QtyReserved: 0}, nil
}
