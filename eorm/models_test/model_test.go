package modelstest

import (
	"eorm"
	"testing"
)

func TestModel(t *testing.T) {

	for _, m := range eorm.ModelRegistry.GetAllModels() {
		t.Log(m)
	}

}
