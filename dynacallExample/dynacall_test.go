package dynacall_test

import (
	"dbx"
	"dynacall"
	_ "dynacall"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "unvs.br.auth"
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
func TestDynaCall(t *testing.T) {
	postData := []interface{}{"Admin", "admin", time.Now()}
	jsonBytes, err := json.Marshal(postData)
	assert.NoError(t, err)
	fmt.Println(string(jsonBytes))

	strJson := string(jsonBytes)
	data := []interface{}{"", "", time.Now()}
	err = json.Unmarshal([]byte(strJson), &data)
	assert.NoError(t, err)
	fmt.Println(data)

	ret, err := dynacall.Call("auth.login@unvs.br.auth", data, struct {
		Tenant   string
		TenantDb *dbx.DBXTenant
	}{
		Tenant:   "testDb",
		TenantDb: createTenantDbx("testDb"),
	})
	assert.Error(t, err)
	fmt.Println(ret)

	list := dynacall.GetAllCaller()
	for _, caller := range list {
		fmt.Println(caller.CallerPath)

	}
}
