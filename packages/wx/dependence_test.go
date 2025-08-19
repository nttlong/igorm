package wx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DependTest struct {
}
type ServerMock struct {
}

func (s *ServerMock) Start(dt *DependTest) error {
	fmt.Println("Http server start")
	return nil
}
func (dp *DependTest) DoRun(s *ServerMock) {
	fmt.Println("DoRun")
}

type AppTest struct {
	DependTest *Depend[DependTest]
	Server     *Depend[ServerMock]
}

func (app *AppTest) New() error {
	app.DependTest.Init(func() (*DependTest, error) {
		return &DependTest{}, nil
	})
	app.Server.Init(func() (*ServerMock, error) {
		return &ServerMock{}, nil
	})
	return nil

}
func TestDependency(t *testing.T) {
	var zapp *AppTest
	err := Start(func(app *AppTest) error {
		zapp = app
		server, err := app.Server.Ins()
		if err != nil {
			return err
		}
		dp, err := app.DependTest.Ins()
		if err != nil {
			return err
		}
		server.Start(dp)
		dp.DoRun(server)
		return nil
	})
	assert.NoError(t, err)
	fmt.Println(zapp)
}
