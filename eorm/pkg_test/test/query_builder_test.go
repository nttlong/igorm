package test

import (
	"eorm"
	"fmt"
	"sync"
	"testing"

	"eorm/migrate"
	"eorm/pkg_test/models"
	_ "eorm/pkg_test/models"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
	db, err := eorm.Open("mssql", msssqlDns)
	assert.NoError(t, err)
	defer db.Close()

	m, err := eorm.NewMigrator(db)
	assert.NoError(t, err)
	err = m.DoMigrates()
	assert.NoError(t, err)

}
func TestCreateObject(t *testing.T) {
	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
	db, err := eorm.Open("mssql", msssqlDns)
	assert.NoError(t, err)

	(&models.User{}).Insert(db)
}

func TestCreateMigrate(t *testing.T) {
	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
	db, err := eorm.Open("mssql", msssqlDns)
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
	m := []migrate.IMigrator{}
	errs := []error{}

	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		db, err := eorm.Open("mssql", msssqlDns)
		assert.NoError(t, err)
		assert.NotEmpty(t, db)
		wg.Add(1) // Quan trọng!

		go func() {
			defer wg.Done()
			msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
			db, err := eorm.Open("mssql", msssqlDns)
			errs = append(errs, err)
			if err == nil {
				rm, err := eorm.NewMigrator(db)
				if err != nil {
					errs = append(errs, err)
				} else {
					m = append(m, rm)
				}

			}
		}()
	}

	wg.Wait() // Đợi cả 2 goroutine xong

	for _, err := range errs {
		if err != nil {
			fmt.Println(err)
		}
		assert.NoError(t, err)
	}
	mcheck := m[0]
	for i := 1; i < len(m); i++ {
		assert.Equal(t, mcheck, m[i])
	}
}

func BenchmarkCreateMigrateTest(b *testing.B) {
	for i := 0; i < b.N; i++ {

		const parallel = 1000
		m := make([]migrate.IMigrator, parallel)
		errs := make([]error, parallel)

		var wg sync.WaitGroup
		wg.Add(parallel)

		for j := 0; j < parallel; j++ {
			j := j // capture đúng j
			go func(index int) {
				defer wg.Done()

				msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
				db, err := eorm.Open("mssql", msssqlDns)
				errs[index] = err
				if err == nil {
					rm, err := eorm.NewMigrator2(db)
					errs[index] = err
					m[index] = rm
				}
			}(j)
		}

		wg.Wait()

		for _, err := range errs {
			assert.NoError(b, err)
		}

		mcheck := m[0]
		for i := 1; i < len(m); i++ {
			assert.Equal(b, mcheck, m[i])
		}
	}
}
