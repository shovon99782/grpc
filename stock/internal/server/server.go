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
