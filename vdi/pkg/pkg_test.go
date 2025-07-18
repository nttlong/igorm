package pkg_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	di "vdi"

	"github.com/stretchr/testify/assert"
)

type TestStructHasFunction struct {
	data string
	Init func(owner *TestStructHasFunction) string
}

func (s *TestStructHasFunction) GetData() string {
	return s.data
}
func TestChangeFunc(t *testing.T) {
	s1 := &TestStructHasFunction{data: "test"}
	s1.Init = func(owner *TestStructHasFunction) string {
		return "test " + owner.data
	}

	// fieldInitOfS1 := reflect.ValueOf(s1).Elem().FieldByName("Init") //lay ham innit ben s1
	s2 := &TestStructHasFunction{data: "test2"} // tao s2
	// gan ham int cua s1 vao s2
	fieldInitOfS1 := reflect.ValueOf(s1).Elem().FieldByName("Init") // hàm Init của s1
	fieldInitOfS2 := reflect.ValueOf(s2).Elem().FieldByName("Init") // hàm Init của s2

	fieldInitOfS2.Set(fieldInitOfS1)
	s2.Init(s2)

}
func TestWriteFileIfExist(t *testing.T) {
	_filePath := `D:\code\go\igorm\unvs-core\unvs\DISCARD`
	if _, err := os.Stat(_filePath); os.IsNotExist(err) {
		fmt.Println("file not exist")
	} else {
		//open file and do something
		err := os.WriteFile(_filePath, ([]byte)("test"), 0644)
		if err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}

		fmt.Printf("Content successfully written to %s\n", _filePath)

		// You can also read the file back to verify the content (optional)
		readContent, err := os.ReadFile(_filePath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		fmt.Printf("Content of %s:\n%s\n", _filePath, string(readContent))
	}
}

func TestDI(t *testing.T) {
	type BaseTestService struct {
		Code di.Singleton[*BaseTestService, string]
	}

	container, err := di.RegisterContainer(func(obj *BaseTestService) error {
		obj.Code.Init = func(owner *BaseTestService) string {
			ret := "code from base"
			return ret
		}
		return nil
	})

	type TestService struct {
		BaseTestService
		Name di.Singleton[*TestService, string]
	}
	_, err = di.RegisterContainer(func(obj *TestService) error {
		obj.Name.Init = func(owner *TestService) string {
			ret := "test"
			return ret
		}
		return nil
	})
	if err != nil {
		assert.NoError(t, err)
	}
	fmt.Println(container.GetInitFun("Name"))
	svc, err := di.Resolve[TestService]()
	assert.NoError(t, err)

	assert.NotEmpty(t, svc)
	value := svc.Code.Get()
	assert.Equal(t, "base", value)
}
