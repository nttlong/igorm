package uorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Model
	UserId   Field
	UserName Field
	Email    Field
}
type TestStruct1 struct {
	name *string
}
type TestStruct2 struct {
	name *string
}

func TestDataTypes(t *testing.T) {
	for i := 0; i < 10; i++ {
		qr := Queryable[User](DB_TYPE_MSSQL, "users")
		assert.Equal(t, "[users].[email] + [users].[user_id]", qr.Email.Add(qr.UserId).String())

		qr2 := qr.As("VVVV").(User) //<-- change table name
		assert.Equal(t, "[VVVV].[email] + [users].[user_id]", qr2.Email.Add(qr.UserId).String())
	}
}
func BenchmarkTestDataTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qr := Queryable[User](DB_TYPE_MSSQL, "users")
		assert.Equal(b, "[users].[email] + [users].[user_id]", qr.Email.Add(qr.UserId).String())

		qr2 := qr.As("VVVV").(User) //<-- change table name
		assert.Equal(b, "[VVVV].[email] + [users].[user_id]", qr2.Email.Add(qr.UserId).String())
	}

}
func TestSelect(t *testing.T) {
	qr := Queryable[User](DB_TYPE_MSSQL, "users")
	assert.Equal(t, "SELECT [users].[email] + [users].[user_id] FROM [users]", qr.Selector(qr.Email.Add(qr.UserId)).String())
}
func BenchmarkTestSelect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//for j := 0; j < 10; j++ {
		qr := Queryable[User](DB_TYPE_MSSQL, "users")

		selector := qr.Selector(qr.Email.Add("qr.UserId"), qr.UserName.Add("qr.UserName"))
		sql := selector.String()
		assert.Equal(b, "SELECT [users].[email] + ?, [users].[user_name] + ? FROM [users]", sql)
		//}
	}

}
