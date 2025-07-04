package orm_test

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestRepository_MSSQL(t *testing.T) {
	for i := 0; i < 10; i++ {
		// mssqlDns := "server=localhost;database=master;user id=sa;password=123456;app name=test"
		// db, err := orm.Open("mssql", mssqlDns)
		// assert.NoError(t, err)
		// defer db.Close()
		// typ := reflect.TypeFor[OrderRepository]()
		// b, err := internal.BuildRepositoryFromType(typ)
		// assert.NoError(t, err)
		// assert.Equal(t, 2, len(b.EntityTypes))
		// assert.Equal(t, "Order", b.EntityTypes[0].Name())
		// assert.Equal(t, "OrderData", b.EntityTypes[1].Name())
		// assert.NotEmpty(t, b.PtrValueOfRepo)
		// repoInstance := b.PtrValueOfRepo.Interface()
		// assert.NotEmpty(t, repoInstance)
		// assert.IsType(t, repoInstance, &OrderRepository{})
		// repo := repoInstance.(*OrderRepository)
		// assert.NotEmpty(t, repo.Base)
		// t.Log(b)
		ret := orm.Repository[OrderRepository]()
		assert.NoError(t, ret.Err)

	}

}
