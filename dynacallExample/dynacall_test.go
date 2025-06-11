package dynacall_test

import (
	"context"
	"dbx"
	"dynacall"
	_ "dynacall"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "unvs.br.auth/roles"
	_ "unvs.br.auth/users"
)

func createCfg() *dbx.Cfg {
	return &dbx.Cfg{
		Driver:   "mssql",
		Host:     "localhost",
		Port:     0,
		User:     "sa",
		Password: "123456",
		SSL:      false,
	}
}
func createDbx() *dbx.DBX {
	ret := dbx.NewDBX(*createCfg())
	return ret
}
func createTenantDbx(tenant string) *dbx.DBXTenant {
	db := createDbx()
	r, e := db.GetTenant(tenant)
	if e != nil {
		panic(e)
	}
	return r
}
func TestShowList(t *testing.T) {
	list := dynacall.GetAllCaller()
	for _, caller := range list {
		fmt.Println(caller.CallerPath)

	}
}
func TestCreateUser(t *testing.T) {
	callPath := "create@unvs.br.auth.users"

	tanentDb := createTenantDbx("testDb")
	tanentDb.Open()
	defer tanentDb.Close()
	for i := 0; i < 1000; i++ {
		username := fmt.Sprintf("user%d", i)
		password := fmt.Sprintf("password%d", i)
		email := fmt.Sprintf("user%d@test.com", i)
		postData := []interface{}{username, password, email}
		start := time.Now()
		ret, err := dynacall.Call(callPath, postData, struct {
			Tenant   string
			TenantDb *dbx.DBXTenant
		}{
			Tenant:   "testDb",
			TenantDb: tanentDb,
		})
		n := time.Since(start).Milliseconds()
		fmt.Printf("time:%d ms\n", n)
		assert.Error(t, err)
		fmt.Println(ret)
	}

}
func configGetJwtSecret() []byte {
	return []byte("super_secret_test_key_for_development_and_testing_only_1234567890ABCDEF")
}
func TestDynaCall(t *testing.T) {
	callPath := "login@unvs.br.auth.users"
	var jwtSecret = configGetJwtSecret()

	tanentDb := createTenantDbx("testDb")
	tanentDb.Open()
	defer tanentDb.Close()
	for i := 0; i < 1000; i++ {
		username := fmt.Sprintf("user%d", i)
		password := fmt.Sprintf("password%d", i)

		postData := []interface{}{username, password}
		start := time.Now()
		ret, err := dynacall.Call(callPath, postData, struct {
			Tenant    string
			TenantDb  *dbx.DBXTenant
			Context   context.Context
			JwtSecret []byte
		}{
			Tenant:    "testDb",
			TenantDb:  tanentDb,
			Context:   context.Background(),
			JwtSecret: jwtSecret,
		})
		n := time.Since(start).Milliseconds()
		fmt.Printf("time:%d ms\n", n)
		t.Log(err)
		t.Log(ret)
		//fmt.Println(ret)
	}

}

type roleInjector struct {
	Tenant      string
	TenantDb    *dbx.DBXTenant
	AccessToken string
	JwtSecret   []byte
}

func (ri *roleInjector) validate() error {
	return nil
}
func TestCreteRole(t *testing.T) {
	callPath := "create@unvs.br.auth.roles"

	tanentDb := createTenantDbx("testDb")
	tanentDb.Open()

	injector := roleInjector{
		Tenant:      "testDb",
		TenantDb:    tanentDb,
		AccessToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDk2MzkzNTgsImlhdCI6MTc0OTYzNTc1OCwicm9sZSI6InVzZXIiLCJzY29wZSI6InJlYWQgd3JpdGUiLCJ1c2VySWQiOiJkNDE4N2YyZC05NzQzLTQwOTgtYWM3MC1mNTcxNzRiNzMyZDIifQ.wQqMrpOR96zx_0hoJlGj4Etk-QmQXF_rSqUzvMTDRYE`,
		JwtSecret:   configGetJwtSecret(),
	}

	defer tanentDb.Close()
	for i := 0; i < 1000; i++ {
		roleName := fmt.Sprintf("role A %d", i)
		roleCode := fmt.Sprintf("A%d", i)
		description := fmt.Sprintf("description %d", i)
		postData := []interface{}{roleCode, roleName, description}
		start := time.Now()
		ret, err := dynacall.Call(callPath, postData, injector)
		n := time.Since(start).Milliseconds()
		fmt.Printf("time:%d ms\n", n)
		assert.Error(t, err)
		fmt.Println(ret)
	}

}
