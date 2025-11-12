package server

import (
	"context"
	"log"

	orderpb "github.com/example/order-service/proto/order"
	// stockpb "github.com/example/stock-service/proto/stock"
)

type orderServer struct {
	orderpb.UnimplementedOrderServiceServer
}

func NewOrderServer() *orderServer {
	return &orderServer{}
}

func (s *orderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Println("CreateOrder called - stub")
	// TODO: validate, call StockService, save to DB, publish event
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
