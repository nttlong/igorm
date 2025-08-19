package wx

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
}
type TestStructWithGlobal struct {
	Cfg *TestConfig
}

func (ts *TestStructWithGlobal) New(Cfg Global[TestConfig], Cfg2 Depend[TestConfig]) error {
	Cfg.Init(func() (*TestConfig, error) {
		return &TestConfig{}, nil
	})
	Cfg2.Init(func() (*TestConfig, error) {

		return &TestConfig{}, nil
	})
	cf2, err := Cfg2.Ins()
	if err != nil {
		return err
	}
	fmt.Println(cf2)
	cf, err := Cfg.Ins()
	if err != nil {
		return err
	}
	ts.Cfg = cf

	return nil
}
func TestCreateGlobal(t *testing.T) {
	mt := GetMethodByName[TestStructWithGlobal]("New")
	receiver := reflect.New(mt.Type.In(0).Elem())
	ret, err := DepenResolvers.RunNewMethodWithReceiver(receiver, *mt)
	assert.Nil(t, err)
	assert.NotNil(t, ret)

	// ok := isTypeDepen(reflect.TypeOf(TestStructWithGlobal{}), map[reflect.Type]bool{})

	// if !ok {
	// 	t.Fatal("not ok")
	// }
	// mt, err := DepenResolvers.FindNewMethod(reflect.TypeOf(TestStructWithGlobal{}))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if mt == nil {
	// 	t.Fatal("mt is nil")
	// }
	// ret, err := DepenResolvers.RunNewMethod(*mt)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if ret == nil {
	// 	t.Fatal("ret is nil")
	// }

}
