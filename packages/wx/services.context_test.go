package wx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestContext1 struct {
	Service
}

func (t *TestContext1) New() error {
	return nil
}

type TestContext2 struct {
}

func TestCheckServiceContext(t *testing.T) {
	ok := serviceHelper.IsTypeHasServiceContext(reflect.TypeFor[TestContext1]())
	assert.True(t, ok)
	ok = serviceHelper.IsTypeHasServiceContext(reflect.TypeFor[TestContext2]())
	assert.False(t, ok)
}
func TestFindNewMethodOfServiceContext(t *testing.T) {
	mt, err := serviceHelper.FindNewMethod(reflect.TypeFor[TestContext1]())
	assert.NoError(t, err)
	assert.NotNil(t, mt)

}
