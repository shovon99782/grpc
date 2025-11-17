package server

import (
	"context"

	"github.com/example/order-service/internal/service"
	orderpb "github.com/example/order-service/proto/order"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	svc *service.OrderService
}

func NewOrderServer(svc *service.OrderService) *OrderServer {
	return &OrderServer{svc: svc}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	return s.svc.CreateOrder(ctx, req)
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.Empty, error) {
	return s.svc.UpdateOrderStatus(ctx, req)
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	return s.svc.GetOrder(ctx, req)
}

func (s *OrderServer) GetOrdersByCustomer(ctx context.Context, req *orderpb.GetOrdersByCustomerRequest) (*orderpb.GetOrdersByCustomerResponse, error) {
	return s.svc.GetOrdersByCustomer(ctx, req)
}
