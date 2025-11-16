package main

import (
	"log"
	"net"

	sql "github.com/example/stock-service/internal/db"
	server "github.com/example/stock-service/internal/server"
	pb "github.com/example/stock-service/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	db := sql.NewMySQLConnection()
	srv := server.NewStockServer(db)
	pb.RegisterStockServiceServer(s, srv)

	log.Println("Stock Service listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
