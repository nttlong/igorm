package depens

import (
	"reflect"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type Test001 struct {
}
type Test002 struct {
}
type Test003 struct {
	T1 *wx.Depend[Test001, Test003]
	T2 *wx.Depend[Test002, Test003]
}
type Args1 struct {
	T3 wx.Depend[Test003, Args1]
}

func (t *Test001) New() error {
	return nil
}
func (t *Test002) New() error {
	return nil
}
func (t *Test003) New(arg *wx.Depend[Args1, Test003], args2 *wx.Global[Args1]) error {
	//_, err := t.T1.Ins()
	// if err != nil {
	// 	return err
	// }

	return nil
}
func TestNewDependTest003(t *testing.T) {
	val, err := wx.Helper.DependNew(reflect.TypeOf(&Test003{}))
	assert.NoError(t, err)
	assert.NotNil(t, val)

}
func TestNewDependTest003Once(t *testing.T) {
	val, err := wx.Helper.DependNewOnce(reflect.TypeOf(&Test003{}))
	assert.NoError(t, err)
	assert.NotNil(t, val)

}
func BenchmarkTestNewDependTest003(t *testing.B) {
	for i := 0; i < t.N; i++ {
		wx.Helper.DependNew(reflect.TypeOf(&Test003{}))

	}

}
func BenchmarkTestNewDependTest003Once(t *testing.B) {
	for i := 0; i < t.N; i++ {

		wx.Helper.DependNewOnce(reflect.TypeOf(&Test003{}))

	}

}
func TestNewDepend(t *testing.T) {
	val, err := wx.NewDepen[Test003]()
	assert.NoError(t, err)
	assert.NotEmpty(t, val)
	t1 := val.T1
	t2 := val.T2
	assert.NotNil(t, t1)
	assert.NotNil(t, t2)
}
func BenchmarkTestNewDepend(t *testing.B) {
	for i := 0; i < t.N; i++ {

		val, err := wx.NewDepen[Test003]()
		assert.NoError(t, err)
		assert.NotEmpty(t, val)
		t1 := val.T1
		t2 := val.T2
		assert.NotNil(t, t1)
		assert.NotNil(t, t2)
	}
}
