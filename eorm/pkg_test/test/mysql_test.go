package test

import (
	"dbv"
	"dbv/tenantDB"
	"fmt"
	"testing"
	"time"

	"dbv/pkg_test/models"
	_ "dbv/pkg_test/models"

	"github.com/stretchr/testify/assert"
)

type HrmRepo struct {
	users *dbv.Repository[models.User]
}

func TestSelectAll(t *testing.T) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := dbv.NewMigrator(db)
	assert.NoError(t, err)
	err = migrator.DoMigrates()
	if err != nil {
		fmt.Println(err)
	}
	rows, err := dbv.SelectAll[models.User](db)
	assert.NoError(t, err)
	assert.Greater(t, len(rows), 151073) //<--lay 151.073 dong

}
func BenchmarkSelectAll(b *testing.B) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(b, err)

	migrator, err := dbv.NewMigrator(db)
	assert.NoError(b, err)
	err = migrator.DoMigrates()
	if err != nil {
		fmt.Println(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := dbv.SelectAll[models.User](db)
		assert.NoError(b, err)
		assert.Greater(b, len(rows), 151000) //<--lay 151.073 dong
	}

}
func BenchmarkSelectAllOldVersion(b *testing.B) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(b, err)

	migrator, err := dbv.NewMigrator(db)
	assert.NoError(b, err)
	err = migrator.DoMigrates()
	if err != nil {
		fmt.Println(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := dbv.SelectAllZeroAllocScanDirect[models.User](db)
		assert.NoError(b, err)
		assert.Greater(b, len(rows), 151000) //<--lay 151.073 dong
	}

}
func TestMySqlKInsertUser(t *testing.T) {

	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := dbv.NewMigrator(db)
	assert.NoError(t, err)
	err = migrator.DoMigrates()
	if err != nil {
		fmt.Println(err)
	}

	assert.NoError(t, err)
	pos := &models.Position{
		Name:  "Manager",
		Level: 1,
		BaseModel: models.BaseModel{
			CreatedAt:   time.Now(),
			Description: dbv.Ptr("test"),
		},
	}
	dep := &models.Department{
		Name: "HRM",
		Code: "HR Department",
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
		},
	}

	err = dbv.Insert(db, dep)
	assert.NoError(t, err)
	err = dbv.Insert(db, pos)
	assert.NoError(t, err)
	user := &models.User{
		Name:       "John",
		Email:      "john@example.com",
		Gender:     "male",
		Birthday:   time.Now(),
		Phone:      "1234567890",
		Address:    "Beijing",
		DeptID:     dep.ID,
		PositionID: pos.ID,
		Username:   nil,
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
		},
	}
	err = dbv.Insert(db, user)
	assert.NoError(t, err)

}
func TestMySqlInsertUserBatch(t *testing.T) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)
	assert.NoError(t, err)
	migrator, err := dbv.NewMigrator(db)
	assert.NoError(t, err)
	err = migrator.DoMigrates()
	assert.NoError(t, err)

	users := make([]*models.User, 0)
	for i := 0; i < 1000; i++ {
		user := &models.User{
			Name:       "John",
			Email:      "john@example.com",
			Gender:     "male",
			Birthday:   time.Now(),
			Phone:      "1234567890",
			Address:    "Beijing",
			DeptID:     1,
			PositionID: 1,
			Username:   nil,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
			},
		}
		users = append(users, user)
	}
	rows, err := dbv.InsertBatch(db, users)
	assert.NoError(t, err)
	assert.Greater(t, rows, int64(0))
}
func BenchmarkMySqlInsertUserBatch(b *testing.B) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)
	assert.NoError(b, err)
	migrator, err := dbv.NewMigrator(db)
	assert.NoError(b, err)
	err = migrator.DoMigrates()
	assert.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		users := make([]*models.User, 0)
		for i := 0; i < 1000; i++ {
			user := &models.User{
				Name:       "John",
				Email:      "john@example.com",
				Gender:     "male",
				Birthday:   time.Now(),
				Phone:      "1234567890",
				Address:    "Beijing",
				DeptID:     1,
				PositionID: 1,
				Username:   nil,
				BaseModel: models.BaseModel{
					CreatedAt: time.Now(),
				},
			}
			users = append(users, user)
		}
		rows, err := dbv.InsertBatch(db, users)
		assert.NoError(b, err)
		assert.Greater(b, rows, int64(0))
	}
}

func BenchmarkMySqlInsertUser(b *testing.B) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)
	assert.NoError(b, err)
	migrator, err := dbv.NewMigrator(db)
	assert.NoError(b, err)
	err = migrator.DoMigrates()
	assert.NoError(b, err)
	tx, err := db.Begin()
	assert.NoError(b, err)
	tx.Db.Begin()
	defer tx.Db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// pos := &models.Position{
		// 	Name:  "Manager",
		// 	Level: 1,
		// 	BaseModel: models.BaseModel{
		// 		CreatedAt:   time.Now(),
		// 		Description: eorm.Ptr("test"),
		// 	},
		// }
		// dep := &models.Department{
		// 	Name: "HRM",
		// 	Code: "HR Department",
		// 	BaseModel: models.BaseModel{
		// 		CreatedAt: time.Now(),
		// 	},
		// }

		// err = eorm.Insert(db, dep)
		// assert.NoError(b, err)
		// err = eorm.Insert(db, pos)
		// assert.NoError(b, err)
		user := &models.User{
			Name:       "John",
			Email:      "john@example.com",
			Gender:     "male",
			Birthday:   time.Now(),
			Phone:      "1234567890",
			Address:    "Beijing",
			DeptID:     1,
			PositionID: 1,
			Username:   nil,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
			},
		}
		err = dbv.InsertWithTx(tx, user)
		assert.NoError(b, err)
	}
	err = tx.Commit()
	assert.NoError(b, err)
}
func BenchmarkMySqlInsertUserRawSql(b *testing.B) {
	sql := "INSERT INTO `users` (`created_at`, `updated_at`, `description`, `name`, `email`, `gender`, `birthday`, `phone`, `address`, `dept_id`, `position_id`, `username`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		user := &models.User{
			Name:       "John",
			Email:      "john@example.com",
			Gender:     "male",
			Birthday:   time.Now(),
			Phone:      "1234567890",
			Address:    "Beijing",
			DeptID:     1,
			PositionID: 1,
			Username:   nil,
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
			},
		}
		r, err := db.Exec(sql, user.CreatedAt, user.UpdatedAt, user.Description, user.Name, user.Email, user.Gender, user.Birthday, user.Phone, user.Address, user.DeptID, user.PositionID, user.Username)
		assert.NoError(b, err)
		rowsAff, err := r.RowsAffected()
		assert.NoError(b, err)
		assert.Greater(b, rowsAff, int64(0))
		id, err := r.LastInsertId()
		assert.NoError(b, err)

		assert.Greater(b, id, int64(0))

	}

}
