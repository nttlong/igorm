package main

import (
	"log"
	"net"

	sv "grpc_container"
	invokerpb "grpc_container/invoker" // 👈 Đặt alias rõ ràng

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Implement service

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	invokerpb.RegisterInvokerServer(s, &sv.Server{})

	// Dòng này phải nằm ở đây
	reflection.Register(s)

	log.Println("gRPC server listening at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
