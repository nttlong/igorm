package depens

import (
	"reflect"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

func TestDepens(t *testing.T) {
	type NoDepens struct {
	}
	type Test struct {
		Fx wx.Depend[Test]
	}
	type Test2 struct {
		Fx wx.Depend[Test]
	}
	ok := wx.Helper.IsGenericDepen(reflect.TypeOf(wx.Depend[Test]{}))
	assert.True(t, ok)
	ok = wx.Helper.IsGenericDepen(reflect.TypeOf(&wx.Depend[Test]{}))
	assert.True(t, ok)
	ok = wx.Helper.IsGenericDepen(reflect.TypeOf(Test{}))
	assert.False(t, ok)
	ok = wx.Helper.IsGenericDepen(reflect.TypeOf(&NoDepens{}))
	assert.False(t, ok)
	fieldIndex := wx.Helper.GetDepen(reflect.TypeOf(wx.Depend[Test]{}))
	assert.Equal(t, [][]int{}, fieldIndex)
	fieldIndex = wx.Helper.GetDepen(reflect.TypeOf(Test{}))
	assert.Equal(t, [][]int{{0}}, fieldIndex)
	fieldIndex = wx.Helper.GetDepen(reflect.TypeOf(&Test{}))
	assert.Equal(t, [][]int{{0}}, fieldIndex)
	fieldIndex = wx.Helper.GetDepen(reflect.TypeOf(Test2{}))
	assert.Equal(t, [][]int{{0}}, fieldIndex)
}
func BenchmarkTestDepens(t *testing.B) {
	type NoDepens struct {
	}
	type Test struct {
		Fx wx.Depend[Test]
	}
	type Test2 struct {
		Fx wx.Depend[Test]
	}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {

		// wx.Helper.IsGenericDepen(reflect.TypeOf(wx.Depend[Test, Test]{}))

		// wx.Helper.IsGenericDepen(reflect.TypeOf(&wx.Depend[Test, Test]{}))

		// wx.Helper.IsGenericDepen(reflect.TypeOf(Test{}))

		// wx.Helper.IsGenericDepen(reflect.TypeOf(&NoDepens{}))

		wx.Helper.GetDepen(reflect.TypeOf(wx.Depend[Test]{}))

		// wx.Helper.GetDepen(reflect.TypeOf(Test{}))

		// wx.Helper.GetDepen(reflect.TypeOf(&Test{}))

		// wx.Helper.GetDepen(reflect.TypeOf(Test2{}))

	}

}

type StructWithNewMethodNotReturnError struct {
}

func (st *StructWithNewMethodNotReturnError) New() {

}

type StructWithNewMethodReturnError struct {
}

func (st *StructWithNewMethodReturnError) New() error {
	return nil
}
func TestGetNewMethodOfStructWithNewMethodNotReturnError(t *testing.T) {
	mt, err := wx.Helper.DependFindNewMethod(reflect.TypeOf(StructWithNewMethodNotReturnError{}))
	assert.Error(t, err)
	assert.Nil(t, mt)

}
func BenchmarkTestGetNewMethodOfStructWithNewMethodNotReturnError(t *testing.B) {
	for i := 0; i < t.N; i++ {
		wx.Helper.DependFindNewMethod(reflect.TypeOf(StructWithNewMethodNotReturnError{}))
		// assert.Error(t, err)
		// assert.Nil(t, mt)
	}

}
func TestGetNewMethodOfStructWithNewMethodReturnError(t *testing.T) {
	mt, err := wx.Helper.DependFindNewMethod(reflect.TypeOf(StructWithNewMethodReturnError{}))
	assert.NoError(t, err)
	assert.NotEmpty(t, mt)

}
func BenchmarkTestGetNewMethodOfStructWithNewMethodReturnError(t *testing.B) {
	for i := 0; i < t.N; i++ {
		wx.Helper.DependFindNewMethod(reflect.TypeOf(StructWithNewMethodReturnError{}))
		// assert.NoError(t, err)
		// assert.NotEmpty(t, mt)
	}

}
