package main

import (
	"fmt"
	"log"
	"net"
	"time"

	caller "vgrpc/caller" // 👈 Đặt alias rõ ràng

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Implement service
type TestService struct {
}
type InputData struct {
	Code string
	Name string
}
type OutputData struct {
	Code        string
	Name        string
	Description string
}

func (svc *TestService) NoArgs() {
	fmt.Println("OK")

}
func (svc *TestService) Run(input InputData) OutputData {
	time.Sleep(time.Second * 5)

	return OutputData{
		Code:        input.Code,
		Name:        input.Name,
		Description: "dasda dasdsa adad ad",
	}

}

func main() {
	caller.AddSingletonService(func() (*TestService, error) {
		return &TestService{}, nil
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	caller.RegisterInvokerServer(s, &caller.Server{})

	// Dòng này phải nằm ở đây
	reflection.Register(s)

	log.Println("gRPC server listening at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
