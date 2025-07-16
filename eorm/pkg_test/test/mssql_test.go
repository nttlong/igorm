package test

import (
	"dbv"
	"dbv/pkg_test/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMssqlInsertDepartment(t *testing.T) {
	msssqlDns := "sqlserver://sa:123456@localhost?database=a002"
	db, err := dbv.Open("sqlserver", msssqlDns)
	assert.NoError(t, err)
	err = db.Ping()

	assert.NoError(t, err)
	m, err := dbv.NewMigrator(db)
	assert.NoError(t, err)
	err = m.DoMigrates()
	assert.NoError(t, err)

	dep := models.Department{
		Name: "HRM",
		Code: "HR Department",
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
		},
	}
	err = db.Insert(&dep)
	assert.NoError(t, err)
}

func TestMssqlInsertDepartmentRaw(t *testing.T) {
	msssqlDns := "sqlserver://sa:123456@localhost?database=a002"
	db, err := dbv.Open("sqlserver", msssqlDns)
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)
	dep := models.Department{
		Name: "HRM",
		Code: "HR Department",
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
		},
	}
	sql := "INSERT INTO [departments] ([name], [code], [parent_id], [created_at], [updated_at], [description]) OUTPUT INSERTED.[id] VALUES (@p1, @p2, @p3, @p4, @p5, @p6)"
	row := db.QueryRow(sql, dep.Name, dep.Code, dep.ParentID, dep.CreatedAt, dep.UpdatedAt, dep.Description)
	assert.NoError(t, err)
	var id int
	err = row.Scan(&id)
	assert.NoError(t, err)
	dep.ID = id
	assert.Equal(t, dep.ID, id)

}
func BenchmarkTestMssqlInsertDepartmentRaw(t *testing.B) {
	msssqlDns := "sqlserver://sa:123456@localhost?database=a002"
	db, err := dbv.Open("sqlserver", msssqlDns)
	assert.NoError(t, err)
	err = db.Ping()

	assert.NoError(t, err)
	m, err := dbv.NewMigrator(db)
	assert.NoError(t, err)
	err = m.DoMigrates()
	assert.NoError(t, err)
	for i := 0; i < t.N; i++ {

		dep := models.Department{
			Name: "HRM",
			Code: "HR Department",
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
			},
		}
		err = db.Insert(&dep)
		assert.NoError(t, err)
	}
}
func BenchmarkTestMssqlInsertDepartment(t *testing.B) {
	msssqlDns := "sqlserver://sa:123456@localhost?database=a002"
	db, err := dbv.Open("sqlserver", msssqlDns)
	assert.NoError(t, err)
	err = db.Ping()

	assert.NoError(t, err)
	m, err := dbv.NewMigrator(db)
	assert.NoError(t, err)
	err = m.DoMigrates()
	assert.NoError(t, err)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		dep := models.Department{
			Code: "HRM",
			Name: "HR Department",
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
			},
		}
		err = db.Insert(&dep)
		assert.NoError(t, err)
	}

}
