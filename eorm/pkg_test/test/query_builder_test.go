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

		paralel := 1
		m := make([]migrate.IMigrator, paralel)
		errs := make([]error, paralel)

		var wg sync.WaitGroup

		for i := 0; i < paralel; i++ {

			wg.Add(1) // Quan trọng!

			go func() {
				defer wg.Done()
				msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
				db, err := eorm.Open("mssql", msssqlDns)
				errs = append(errs, err)
				if err == nil {
					rm, err := eorm.NewMigrator(db)
					errs[i] = err
					m[i] = rm

				}
			}()
		}

		wg.Wait() // Đợi cả 2 goroutine xong

		for _, err := range errs {
			if err != nil {
				fmt.Println(err)
			}
			assert.NoError(b, err)
		}
		mcheck := m[0]
		for i := 1; i < len(m); i++ {
			assert.Equal(b, mcheck, m[i])
		}
	}
}
