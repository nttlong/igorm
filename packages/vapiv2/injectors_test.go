package vapi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FileUstils struct {
}
type DbContext struct {
}

func (fx *FileUstils) SaveFile() {

}

type FileService struct {
	FileUtils Singleton[FileUstils, FileService]
	Db        Transient[DbContext, FileService]
}

func TestInjectorSingleton(t *testing.T) {
	RegisterService(func(svc *FileService) error {

		svc.FileUtils.Init(func() (*FileUstils, error) {
			fmt.Println(svc)
			return &FileUstils{}, nil
		})
		svc.Db.Init(func() (*DbContext, error) {
			fmt.Println(svc)
			return &DbContext{}, nil
		})
		return nil
	})
	fileService, err := Service[FileService]()
	if err != nil {
		t.Error(err)
	}
	fileService.FileUtils.GetInstance().SaveFile()
	fileUtisl := fileService.FileUtils.GetInstance()
	db, err := fileService.Db.GetInstance()
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NoError(t, err)
	assert.NotNil(t, fileUtisl)
}
func BenchmarkInjectorSingleton(t *testing.B) {
	for i := 0; i < t.N; i++ {
		RegisterService(func(svc *FileService) error {

			svc.FileUtils.Init(func() (*FileUstils, error) {
				fmt.Println("init FileUstils")
				return &FileUstils{}, nil
			})
			svc.Db.Init(func() (*DbContext, error) {
				fmt.Println("init db")
				return &DbContext{}, nil

			})
			return nil
		})
		fileService, err := Service[FileService]()
		if err != nil {
			t.Error(err)
		}
		fileService.FileUtils.GetInstance().SaveFile()
		fileUtisl := fileService.FileUtils.GetInstance()
		db, err := fileService.Db.GetInstance()
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.NoError(t, err)
		assert.NotNil(t, fileUtisl)
	}

}
func TestIsSingletonType(t *testing.T) {
	assert.True(t, serviceUtils.IsSingletonType(reflect.TypeOf(Singleton[any, any]{})))
	assert.False(t, serviceUtils.IsSingletonType(reflect.TypeOf(FileService{})))

}
