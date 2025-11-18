package main

import (
	"log"
	"net"

	"github.com/example/stock-service/config"
	sql "github.com/example/stock-service/internal/db"
	rabbitmq "github.com/example/stock-service/internal/rabbitmq"
	server "github.com/example/stock-service/internal/server"
	service "github.com/example/stock-service/internal/service"
	pb "github.com/example/stock-service/proto"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	db := sql.NewMySQLConnection()
	service := service.NewStockService(db)

	rabbit := rabbitmq.NewRabbitMQ(cfg.RabbitUrl)
	ch, err := rabbit.Conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open channel: %v", err)
	}
	cancelConsumer := rabbitmq.NewOrderCancelConsumer(ch, service)

	err = cancelConsumer.Start()
	if err != nil {
		log.Fatalf("❌ Failed to start cancel consumer: %v", err)
	}
	log.Println("📡 OrderCancelConsumer started")

	srv := server.NewStockServer(service)
	pb.RegisterStockServiceServer(s, srv)

	log.Println("Stock Service listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
