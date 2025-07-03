package orm_test

import (
	"testing"
	orm "unvs-orm"
)

type Order struct {
	*orm.Model[Order]
	OrderId   orm.NumberField[uint64] `db:"primaryKey;autoIncrement"`
	Version   orm.NumberField[int]    `db:"primaryKey"`
	Note      orm.TextField           `db:"length(200)"`
	CreatedAt orm.DateTimeField
	UpdatedAt orm.DateTimeField `db:"null"`
	CreatedBy orm.TextField     `db:"length(100)"`
	UpdatedBy orm.TextField     `db:"length(100);null"`
}
type OrderItem struct {
	*orm.Model[OrderItem]
	Id        orm.NumberField[uint64] `db:"primaryKey;autoIncrement"`
	OrderId   orm.NumberField[uint64] `db:"index(order_ref_idx)"`
	Version   orm.NumberField[int]    `db:"index(order_ref_idx)"`
	Product   orm.TextField           `db:"length(100)"`
	Quantity  orm.NumberField[int]
	CreatedAt orm.DateTimeField
	UpdatedAt orm.DateTimeField `db:"null"`
	CreatedBy orm.TextField     `db:"length(100)"`
	UpdatedBy orm.TextField     `db:"length(100);null"`
}
type OrderRepository struct {
	*orm.TenantDb
	Orders     *Order
	OrderItems *OrderItem
}

func (r *OrderRepository) Init() {
	r.NewRelationship().
		From(r.Orders.OrderId, r.Orders.Version).
		To(r.OrderItems.OrderId, r.OrderItems.Version)
}
func TestGenerateSqlMigrate(t *testing.T) {
}
