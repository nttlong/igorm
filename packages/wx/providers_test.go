package wx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type IService interface {
	GetName() string
}

type Service1 struct {
	Provider[Service1, IService]
}

func (s *Service1) GetName() string {
	return "Service1"
}

func TestRegisterContainer(t *testing.T) {
	(&Service1{}).Register(func(svc *Service1) (IService, error) {
		return svc, nil
	})
	svc, err := (&Service1{}).New()
	assert.NoError(t, err)
	assert.Equal(t, svc.GetName(), "Service1")
}
