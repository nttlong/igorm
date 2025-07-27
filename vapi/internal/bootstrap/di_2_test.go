package bootstrap

import (
	"testing"
	"vdi"

	"github.com/stretchr/testify/assert"
)

func TestContaierTest(t *testing.T) {
	c := (&ContaierTest{}).New(func(owner *ContaierTest) error {
		owner.Number.SetFunc(func(owner *ContaierTest, args ...any) struct{ Number int } {
			val := args[0].(int)

			return struct{ Number int }{val + 5}
		})

		return nil

	})
	for j := 0; j < 10; j++ {
		v := c.Number.Get(j)
		assert.Equal(t, v.Number, j+5)
	}

}

type ContaierTest struct {
	*vdi.Container[ContaierTest]
	Number vdi.Transient[ContaierTest, struct{ Number int }]
}

func BenchmarkTestContaierTest(b *testing.B) {
	c := (&ContaierTest{}).New(func(owner *ContaierTest) error {
		owner.Number.SetFunc(func(owner *ContaierTest, args ...any) struct{ Number int } {
			val := args[0].(int)
			val2 := args[1].(int)

			return struct{ Number int }{val * val2}
		})

		return nil

	})
	for i := 0; i < b.N; i++ {

		for j := 0; j < 10; j++ {
			v := c.Number.Get(j, i)
			assert.Equal(b, i*j, v.Number)
		}
	}

}
