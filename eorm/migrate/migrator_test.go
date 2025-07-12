package migrate

import (
	"eorm/tenantDB"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Model struct {
	Entity
}
type User struct {
	Model `db:"table:users"`

	ID        int       `db:"pk;auto"`                         // primary key, auto increment
	Name      string    `db:"pk;size:50"`                      // mapped column name, varchar(50)
	Email     string    `db:"uk:test_email;size:120"`          // unique constraint named "test_email"
	Profile   *string   `db:"size:255"`                        // nullable string
	CreatedAt time.Time `db:"default:now;type:datetime"`       // default timestamp
	Price     float64   `db:"type:decimal(10,2);column:price"` // custom type and column name
}

func init() {
	ModelRegistry.Add(&User{})
}

func TestMigrator(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := tenantDB.Open("mssql", sqlServerDns)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	migrator, err := NewMigrator(db)
	assert.NoError(t, err)
	tables, err := migrator.GetSqlCreateTable(reflect.TypeOf(User{}))
	assert.NoError(t, err)

	fmt.Print(tables)
	assert.NotEmpty(t, tables)

	// TODO: implement test cases
}
