package invoker

import (
	"context"
	"log"
	"testing"
	"time"

	invokerpb "grpc_container/invoker" // Import lại package đã được tạo từ .proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	serverAddress = "localhost:50051"
)

// TestInvokerClient là hàm test tích hợp
func TestInvokerClient(t *testing.T) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Không thể kết nối đến server gRPC: %v", err)
	}
	defer conn.Close() // Đóng kết nối khi hàm main kết thúc

	// Tạo một client gRPC
	client := invokerpb.NewInvokerClient(conn)

	// Tạo context với timeout để tránh treo chương trình
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tạo dữ liệu cho yêu cầu InvokeRequest
	paramsMap := map[string]interface{}{
		"arg1": "client_value_1",
		"arg2": 456,
	}
	paramsStruct, err := structpb.NewStruct(paramsMap)
	if err != nil {
		log.Fatalf("Không thể tạo structpb.Struct: %v", err)
	}

	// Tạo yêu cầu InvokeRequest
	req := &invokerpb.InvokeRequest{
		PackagePath:  "go_client/package",
		FunctionName: "fail",
		Params:       paramsStruct,
	}

	// Gọi phương thức Invoke trên server
	log.Printf("Đang gọi phương thức Invoke trên server...")
	reply, err := client.Invoke(ctx, req)
	if err != nil {
		panic(err)
	}

	// Xử lý và in ra kết quả
	log.Printf("Gọi thành công! Server trả về: %v", reply.Result)
	// Để lấy giá trị cụ thể từ structpb.Struct
	message := reply.Result.Fields["message"].GetStringValue()
	log.Printf("Giá trị message từ server: %s", message)
}
