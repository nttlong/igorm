package test

// import (
// 	"vdb"
// 	"fmt"
// 	"sync"
// 	"testing"

// 	"vdb/migrate"
// 	"vdb/pkg_test/models"
// 	_ "vdb/pkg_test/models"

// 	"github.com/stretchr/testify/assert"
// )

// func TestForeignKey(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 	db, err := vdb.Open("mssql", msssqlDns)
// 	assert.NoError(t, err)
// 	defer db.Close()
// 	m, err := vdb.NewMigrator(db)
// 	assert.NoError(t, err)

// 	m.DoMigrates()
// }
// func TestQueryBuilder(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 	db, err := vdb.Open("mssql", msssqlDns)
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	m, err := vdb.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)

// }
// func TestCreateObject(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 	db, err := vdb.Open("mssql", msssqlDns)
// 	assert.NoError(t, err)

// 	(&models.User{}).Insert(db)
// }

// func TestCreateMigrate(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 	db, err := vdb.Open("mssql", msssqlDns)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, db)
// 	m := []migrate.IMigrator{}
// 	errs := []error{}

// 	var wg sync.WaitGroup

// 	for i := 0; i < 1; i++ {
// 		db, err := vdb.Open("mssql", msssqlDns)
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, db)
// 		wg.Add(1) // Quan trọng!

// 		go func() {
// 			defer wg.Done()
// 			msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 			db, err := vdb.Open("mssql", msssqlDns)
// 			errs = append(errs, err)
// 			if err == nil {
// 				rm, err := vdb.NewMigrator(db)
// 				if err != nil {
// 					errs = append(errs, err)
// 				} else {
// 					m = append(m, rm)
// 				}

// 			}
// 		}()
// 	}

// 	wg.Wait() // Đợi cả 2 goroutine xong

// 	for _, err := range errs {
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		assert.NoError(t, err)
// 	}
// 	mcheck := m[0]
// 	for i := 1; i < len(m); i++ {
// 		assert.Equal(t, mcheck, m[i])
// 	}
// }

// func BenchmarkCreateMigrateTest(b *testing.B) {
// 	for i := 0; i < b.N; i++ {

// 		const parallel = 1000
// 		m := make([]migrate.IMigrator, parallel)
// 		errs := make([]error, parallel)

// 		var wg sync.WaitGroup
// 		wg.Add(parallel)

// 		for j := 0; j < parallel; j++ {
// 			j := j // capture đúng j
// 			go func(index int) {
// 				defer wg.Done()

// 				msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 				db, err := vdb.Open("mssql", msssqlDns)
// 				errs[index] = err
// 				if err == nil {
// 					rm, err := vdb.NewMigrator(db)
// 					if err == nil {
// 						err = rm.DoMigrates()
// 					}
// 					errs[index] = err
// 					m[index] = rm
// 				}
// 			}(j)
// 		}

// 		wg.Wait()

// 		for _, err := range errs {
// 			assert.NoError(b, err)
// 		}

// 		mcheck := m[0]
// 		for i := 1; i < len(m); i++ {
// 			assert.Equal(b, mcheck, m[i])
// 		}
// 	}
// }
// func TestInsertUser(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 	db, err := vdb.Open("mssql", msssqlDns)
// 	assert.NoError(t, err)
// 	defer db.Close()
// 	for i := 19607; i < 19607+10000; i++ {
// 		data, err := vdb.Repo[models.User]().Insert(db, &models.User{
// 			Name:     "test" + fmt.Sprintf("%d", i),
// 			Email:    "test" + fmt.Sprintf("%d", i) + "@gmail.com",
// 			Username: vdb.Ptr("test" + fmt.Sprintf("%d", i)),
// 		})
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, data)
// 	}

// }
// func TestInsertUserWitTx(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 	db, err := vdb.Open("mssql", msssqlDns)
// 	assert.NoError(t, err)
// 	defer db.Close()
// 	tx, err := db.Begin()
// 	assert.NoError(t, err)
// 	for i := 19607; i < 19607+10000; i++ {
// 		data, err := vdb.Repo[models.User]().InsertWithTx(tx, &models.User{
// 			Name:     "test" + fmt.Sprintf("%d", i),
// 			Email:    "test" + fmt.Sprintf("%d", i) + "@gmail.com",
// 			Username: vdb.Ptr("test" + fmt.Sprintf("%d", i)),
// 		})
// 		if err != nil {

// 			assert.Equal(t, "code=ERR0001, duplicate: duplicate cols username tables users, entity fields Username", err.Error())
// 		} else {
// 			assert.NoError(t, err)
// 		}
// 		assert.NotEmpty(t, data)
// 	}
// 	err = tx.Commit()
// 	assert.NoError(t, err)
// }
// func BenchmarkInsertUser(b *testing.B) {
// 	msssqlDns := "sqlserver://sa:123456@localhost:1433?database=a001"
// 	db, err := vdb.Open("mssql", msssqlDns)
// 	assert.NoError(b, err)
// 	defer db.Close()
// 	for i := 0; i < b.N; i++ {

// 		data, err := vdb.Repo[models.User]().InsertContext(b.Context(), db, &models.User{ //<-- inser to database
// 			Name:     "test" + fmt.Sprintf("%d", i+50000),
// 			Email:    "test" + fmt.Sprintf("%d", i+50000) + "@gmail.com",
// 			Username: vdb.Ptr("test" + fmt.Sprintf("%d", i+30000)),
// 		})
// 		if err != nil {

// 			assert.Equal(b, "code=ERR0001, duplicate: duplicate cols username tables users, entity fields Username", err.Error())
// 		} else {
// 			assert.NoError(b, err)
// 		}

// 		assert.NotEmpty(b, data)
// 	}
// }
