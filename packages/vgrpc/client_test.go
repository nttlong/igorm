package vgrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Userinfo struct {
	Name string
	Age  int
}

func TestCaller(t *testing.T) {
	caller := NewClientCaller("localhost", "50051", 15)
	caller.Connect()
	defer caller.Disconnect()
	// ret, err := caller.Call("main.TestService.NoArgs", nil)
	ret, err := caller.Call("main.TestService.Run", Userinfo{
		Name: "test",
		Age:  1,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)

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

func BenchmarkCaller(b *testing.B) {
	caller := NewClientCaller("localhost", "50051", 15)
	caller.Connect()
	defer caller.Disconnect()
	for i := 0; i < b.N; i++ {

		ret, err := caller.Call("main.TestService.Run", InputData{
			Code: "A001",
			Name: "test",
		})
		assert.NoError(b, err)
		assert.NotEmpty(b, ret)
	}
}
