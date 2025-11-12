package main

import (
    "context"
    pb "github.com/example/order-service/proto"
    "log"
)

type orderServer struct {
    pb.UnimplementedOrderServiceServer
}

func NewOrderServer() *orderServer {
    return &orderServer{}
}

func (s *orderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    log.Println("CreateOrder called - stub")
    // TODO: validate, call StockService, save to DB, publish event
    return &pb.CreateOrderResponse{OrderId: "stub-order-id", Status: "CREATED"}, nil
}

func (s *orderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Empty, error) {
    log.Println("UpdateOrderStatus called - stub")
    return &pb.Empty{}, nil
}

func (s *orderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
    log.Println("GetOrder called - stub")
    return &pb.GetOrderResponse{OrderId: req.OrderId, CustomerId: "cust-stub", Status: "CREATED"}, nil
}

func (s *orderServer) GetOrdersByCustomer(ctx context.Context, req *pb.GetOrdersByCustomerRequest) (*pb.GetOrdersByCustomerResponse, error) {
    log.Println("GetOrdersByCustomer called - stub")
    return &pb.GetOrdersByCustomerResponse{}, nil
}
