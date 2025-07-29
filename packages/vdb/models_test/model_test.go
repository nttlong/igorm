package modelstest

import (
	"testing"
	"vdb"
)

func TestModel(t *testing.T) {

	for _, m := range vdb.ModelRegistry.GetAllModels() {
		t.Log(m)
	}

}
