package services

import (
	"testing"
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
	assert.NoError(t, err)

}
