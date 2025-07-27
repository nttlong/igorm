package vdi

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestService struct {
	BaseService
}

func TestTransisent(t *testing.T) {
	container := TestService{}
	container2 := TestService{}

	TransisentRegister(&container, func() (*int64, error) {
		val := TransisentGet[int64](&container2)
		if container2.Err != nil {
			return nil, container2.Err
		}
		val += 1
		return &val, nil

	})
	TransisentRegister(&container2, func() (*int64, error) {
		ret := time.Now().Unix()
		return &ret, nil

	})
	fmt.Println(TransisentGet[int64](&container))

}
func BenchmarkTestTransisent(b *testing.B) {
	container := TestService{}
	container2 := TestService{}

	TransisentRegister(&container, func() (*int, error) {
		val := TransisentGet[int](&container2)
		if container2.Err != nil {
			return nil, container2.Err
		}
		val += 1
		return &val, nil

	})

	for i := 0; i < b.N; i++ {
		for j := 0; j < 5; j++ {
			TransisentRegister(&container2, func() (*int, error) {
				fmt.Println(j)
				ret := j
				return &ret, nil

			})
			val := TransisentGet[int](&container)
			assert.Equal(b, int(i+1), val)
		}

	}

}
