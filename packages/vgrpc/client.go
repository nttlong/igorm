package vgrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// <-- Sửa lại dòng này

	//"github.com/golang/protobuf/proto" //<-- sau khi chay lenh
	/*
		go mod edit -replace github.com/golang/protobuf=google.golang.org/protobuf@v1.27.1
		go mod tidy
		 cho nay bi loi nhu sau
		 could not import github.com/golang/protobuf/proto (no required module provides package "github.com/golang/protobuf/proto")comp
	*/
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import package được tạo từ proto của bạn
	"vgrpc/caller"
)

type ClientCaller struct {
	conn    *grpc.ClientConn
	client  caller.InvokerClient
	Host    string
	Port    string
	Timeout int
}

func NewClientCaller(host, port string, timeOut int) *ClientCaller {
	return &ClientCaller{
		Host:    host,
		Port:    port,
		Timeout: timeOut,
	}
}
func (c *ClientCaller) Connect() error {
	var err error
	c.conn, err = grpc.Dial(c.Host+":"+c.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := caller.NewInvokerClient(c.conn)
	c.client = client
	return nil
}
func (c *ClientCaller) Disconnect() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// packInputs converts a slice of Go interfaces into a slice of Any messages.

func (c *ClientCaller) Call(methodPath string, input interface{}) (interface{}, error) {
	if input != nil {
		jsonInput, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}

		request := &caller.InvokeRequest{
			MethodPath: methodPath, // Thay thế bằng đường dẫn hàm thực tế
			InputJson:  jsonInput,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(c.Timeout))
		defer cancel()
		response, err := c.client.Invoke(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("Can not call Invoke: %w", err)
		}

		if response.GetErrorMessage() != "" {
			return nil, fmt.Errorf("Server error: %s", response.GetErrorMessage())
		}

		return input, nil
	} else {

		request := &caller.InvokeRequest{
			MethodPath: methodPath, // Thay thế bằng đường dẫn hàm thực tế
			InputJson:  []byte{},
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(c.Timeout))
		defer cancel()
		response, err := c.client.Invoke(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("Can not call Invoke: %w", err)
		}

		if response.GetErrorMessage() != "" {
			return nil, fmt.Errorf("Server error: %s", response.GetErrorMessage())
		}

		return input, nil
	}

}
