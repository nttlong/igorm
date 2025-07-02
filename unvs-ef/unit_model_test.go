package unvsef

import (
	"database/sql"
	"testing"
)

type Order struct {
	Entity[Order]
	OrderId   FieldNumber[uint64] `db:"primaryKey;autoIncrement"`
	Version   FieldNumber[int]    `db:"primaryKey"`
	Note      FieldString         `db:"length(200)"`
	CreatedAt FieldDateTime
	UpdatedAt *FieldDateTime
	CreatedBy FieldString  `db:"length(100)"`
	UpdatedBy *FieldString `db:"length(100)"`
}
type OrderItem struct {
	Entity[OrderItem]
	Id        FieldNumber[uint64] `db:"primaryKey;autoIncrement"`
	OrderId   FieldNumber[uint64] `db:"index(order_ref_idx)"`
	Version   FieldNumber[int]    `db:"index(order_ref_idx)"`
	Product   FieldString         `db:"length(100)"`
	Quantity  FieldNumber[int]
	CreatedAt FieldDateTime
	UpdatedAt *FieldDateTime
	CreatedBy FieldString  `db:"length(100)"`
	UpdatedBy *FieldString `db:"length(100)"`
}
type OrderRepository struct {
	*TenantDb
	Orders     *Order
	OrderItems *OrderItem
}

func (r *OrderRepository) Init() {
	r.NewRelationship().
		From(r.Orders.OrderId, r.Orders.Version).
		To(r.OrderItems.OrderId, r.OrderItems.Version)
}
func TestBinaryField(t *testing.T) {

	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := Repo[OrderRepository](db, true) // create repos
	expr := repo.Orders.OrderId.Gt(1000).And(*repo.Orders.Version.Lt(1000))
	sql, args := expr.ToSqlExpr(repo.Dialect)

	t.Log(sql, args)
	expr = expr.Or(repo.Orders.Note.Eq("test"))
	sql, args = expr.ToSqlExpr(repo.Dialect)

	t.Log(sql, args)
}
