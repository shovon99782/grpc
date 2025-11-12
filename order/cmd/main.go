package main

import (
	"context"
	"log"
	"net"

	server "github.com/example/order-service/internal/server"
	pb "github.com/example/order-service/proto/order"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// register server (implementation in internal/server)
	srv := server.NewOrderServer()
	pb.RegisterOrderServiceServer(s, srv)

	log.Println("Order Service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	_ = context.Background()
}
