package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    pb "github.com/example/stock-service/proto"
)

func main() {
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    srv := NewStockServer()
    pb.RegisterStockServiceServer(s, srv)

    log.Println("Stock Service listening on :50052")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
