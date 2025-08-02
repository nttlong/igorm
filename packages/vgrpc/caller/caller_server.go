package caller

import (
	context "context"
)

type Server struct {
	UnimplementedInvokerServer
}
type Userinfo struct {
	Name string
	Age  int
}

func (s *Server) Invoke(ctx context.Context, req *InvokeRequest) (*InvokeResponse, error) {
	ret, err := Call(req.MethodPath, req.InputJson)
	if err != nil {
		return nil, err
	}
	return &InvokeResponse{
		Result: ret,
	}, nil
}
