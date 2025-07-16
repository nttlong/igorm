package modelstest

import (
	"dbv"
	"testing"
)

func TestModel(t *testing.T) {

	for _, m := range dbv.ModelRegistry.GetAllModels() {
		t.Log(m)
	}

}
