package main

import (
	"context"
	"log"
	"net"

	"github.com/example/order-service/config"
	sql "github.com/example/order-service/internal/db"
	"github.com/example/order-service/internal/rabbitmq"
	server "github.com/example/order-service/internal/server"
	service "github.com/example/order-service/internal/service"
	pb "github.com/example/order-service/proto/order"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// register server (implementation in internal/server)
	db := sql.NewMySQLConnection()
	rabbit := rabbitmq.NewRabbitMQ(cfg.RabbitUrl)
	orderService := service.NewOrderService(db, rabbit)
	// srv := server.NewOrderServer(db, rabbit)
	orderServer := server.NewOrderServer(orderService)
	pb.RegisterOrderServiceServer(grpcServer, orderServer)

	log.Println("Order Service listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	_ = context.Background()
}
