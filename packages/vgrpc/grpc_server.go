package vgrpc

import (
	"context"
	"fmt"

	invokerpb "vgrpc/invoker" // 👈 Đặt alias rõ ràng

	"google.golang.org/protobuf/types/known/structpb"
)

type Server struct {
	invokerpb.UnimplementedInvokerServer
}

func (s *Server) Invoke(ctx context.Context, req *invokerpb.InvokeRequest) (*invokerpb.InvokeReply, error) {
	// Demo: In ra tên hàm được gọi và params
	fmt.Println("C 1 loi goi ham tu client:", req.FunctionName)
	fmt.Println("Params:", req.Params)

	// Ví dụ: giả lập lỗi khi gọi function không tồn tại
	if req.FunctionName == "fail" {
		return nil, fmt.Errorf("function %s not found", req.FunctionName)
	}

	// Trả về reply mẫu
	replyData, _ := structpb.NewStruct(map[string]interface{}{
		"message": "Called " + req.FunctionName,
	})

	return &invokerpb.InvokeReply{
		Result: replyData,
	}, nil
}
