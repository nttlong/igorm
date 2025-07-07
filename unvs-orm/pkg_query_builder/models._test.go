package pkgquerybuilder

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
type Customer struct {
	*orm.Model[Customer]
	CustomerId orm.NumberField[uint64] `db:"primaryKey;autoIncrement"`
	Code       orm.TextField           `db:"length(50);unique"` // Mã khách hàng
	Name       orm.TextField           `db:"length(200)"`       // Tên khách
	Email      orm.TextField           `db:"length(100);null"`
	Phone      orm.TextField           `db:"length(20);null"`
	Address    orm.TextField           `db:"length(300);null"`

	Note      orm.TextField `db:"length(200);null"`
	CreatedAt orm.DateTimeField
	UpdatedAt orm.DateTimeField `db:"null"`
	CreatedBy orm.TextField     `db:"length(100)"`
	UpdatedBy orm.TextField     `db:"length(100);null"`
}
type Invoice struct {
	*orm.Model[Invoice]
	InvoiceId       orm.NumberField[uint64] `db:"primaryKey;autoIncrement"`
	OrderId         orm.NumberField[uint64] `db:"index(order_invoice_idx)"` // Nếu liên kết tới Order
	Version         orm.NumberField[int]    `db:"index(order_invoice_idx)"`
	CustomerId      orm.NumberField[uint64] `db:"index(customer_invoice_idx)"`
	PaymentMethodId orm.NumberField[uint64] `db:"index(payment_method_idx)"` // Mới thêm
	InvoiceDate     orm.DateTimeField
	TotalAmount     orm.NumberField[float64]
	Name            orm.TextField `db:"length(200)"`

	Note      orm.TextField `db:"length(200);null"`
	CreatedAt orm.DateTimeField
	UpdatedAt orm.DateTimeField `db:"null"`
	CreatedBy orm.TextField     `db:"length(100)"`
	UpdatedBy orm.TextField     `db:"length(100);null"`
}

type Item struct {
	*orm.Model[Item]
	ItemId orm.NumberField[uint64]  `db:"primaryKey;autoIncrement"`
	Code   orm.TextField            `db:"length(50);unique"`
	Name   orm.TextField            `db:"length(200)"`
	Unit   orm.TextField            `db:"length(20)"` // Đơn vị tính
	Price  orm.NumberField[float64] // Giá mặc định
	Note   orm.TextField            `db:"length(200);null"`

	CreatedAt orm.DateTimeField
	UpdatedAt orm.DateTimeField `db:"null"`
	CreatedBy orm.TextField     `db:"length(100)"`
	UpdatedBy orm.TextField     `db:"length(100);null"`
}
type InvoiceDetail struct {
	*orm.Model[InvoiceDetail]
	Id orm.NumberField[uint64] `db:"primaryKey;autoIncrement"`

	InvoiceId orm.NumberField[uint64] `db:"index(invoice_detail_idx)"`
	Version   orm.NumberField[int]    `db:"index(invoice_detail_idx)"`

	ItemId orm.NumberField[uint64]  `db:"index"`
	Amount orm.NumberField[int]     // Số lượng
	Price  orm.NumberField[float64] // Đơn giá
	Note   orm.TextField            `db:"length(200);null"`

	CreatedAt orm.DateTimeField
	UpdatedAt orm.DateTimeField `db:"null"`
	CreatedBy orm.TextField     `db:"length(100)"`
	UpdatedBy orm.TextField     `db:"length(100);null"`
}
type OrderRepository struct {
	*orm.Base
	Orders         *Order
	OrderItems     *OrderItem
	Invoices       *Invoice
	Customers      *Customer
	Items          *Item
	InvoiceDetails *InvoiceDetail
	PaymentMethods *PaymentMethod
}
type PaymentMethod struct {
	*orm.Model[PaymentMethod]
	PaymentMethodId orm.NumberField[uint64] `db:"primaryKey;autoIncrement"`
	Code            orm.TextField           `db:"length(50);unique"` // Ví dụ: "CASH", "BANK", "CARD"
	Name            orm.TextField           `db:"length(100)"`       // Ví dụ: "Tiền mặt", "Chuyển khoản"

	Note      orm.TextField `db:"length(200);null"`
	CreatedAt orm.DateTimeField
	UpdatedAt orm.DateTimeField `db:"null"`
	CreatedBy orm.TextField     `db:"length(100)"`
	UpdatedBy orm.TextField     `db:"length(100);null"`
}

func (r *OrderRepository) Init() {
	r.NewRelationship().
		From(r.Orders.OrderId, r.Orders.Version).
		To(r.OrderItems.OrderId, r.OrderItems.Version)
}
func TestGenerateSqlMigrate(t *testing.T) {
}
