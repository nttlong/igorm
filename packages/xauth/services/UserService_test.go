package services

import (
	"testing"
	"vdb"
	"wx"
	dbmodels "xauth/dbModels"

	"github.com/stretchr/testify/assert"
)

func TestUserserviceCreateUser(t *testing.T) {
	db, err := wx.NewGlobal[DbService]()
	assert.NoError(t, err)
	assert.NotNil(t, db)

	ret, err := wx.NewDepend[UserService]()
	assert.NoError(t, err)
	tenantDb, err := db.GetTenantDb("test001")
	assert.NoError(t, err)
	user := &dbmodels.Users{
		Username: "test001",
	}
	err = ret.CreateUser(tenantDb, user, "test001")
	if err != nil {
		var vErr *vdb.DialectError
		if !assert.ErrorAs(t, err, &vErr) {
			t.FailNow()
		}
		assert.Equal(t, vdb.DIALECT_DB_ERROR_TYPE_DUPLICATE, vErr.ErrorType)

	}

}
func BenchmarkUserserviceCreateUser(b *testing.B) {
	db, err := wx.NewGlobal[DbService]()
	assert.NoError(b, err)
	assert.NotNil(b, db)
	assert.NoError(b, err)
	tenantDb, err := db.GetTenantDb("test001")
	assert.NoError(b, err)
	ret, err := wx.NewDepend[UserService]()
	assert.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		user := &dbmodels.Users{
			Username: "test001",
		}
		err = ret.CreateUser(tenantDb, user, "test001")
		if err != nil {
			var vErr *vdb.DialectError
			if !assert.ErrorAs(b, err, &vErr) {
				b.FailNow()
			}
			assert.Equal(b, vdb.DIALECT_DB_ERROR_TYPE_DUPLICATE, vErr.ErrorType)

		}
	}
}
