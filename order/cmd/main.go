package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "github.com/example/order-service/proto"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    // register server (implementation in internal/server)
    srv := NewOrderServer()
    pb.RegisterOrderServiceServer(s, srv)

    log.Println("Order Service listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }

    _ = context.Background()
}
