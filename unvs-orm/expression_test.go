package orm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Order struct {
	*Model[Order]
	OrderId   NumberField[uint64] `db:"primaryKey;autoIncrement"`
	Version   NumberField[int]    `db:"primaryKey"`
	Note      TextField           `db:"length(200)"`
	CreatedAt DateTimeField
	UpdatedAt DateTimeField `db:"null"`
	CreatedBy TextField     `db:"length(100)"`
	UpdatedBy TextField     `db:"length(100);null"`
}
type OrderItem struct {
	*Model[OrderItem]
	Id        NumberField[uint64] `db:"primaryKey;autoIncrement"`
	OrderId   NumberField[uint64] `db:"index(order_ref_idx)"`
	Version   NumberField[int]    `db:"index(order_ref_idx)"`
	Product   TextField           `db:"length(100)"`
	Quantity  NumberField[int]
	CreatedAt DateTimeField
	UpdatedAt DateTimeField `db:"null"`
	CreatedBy TextField     `db:"length(100)"`
	UpdatedBy TextField     `db:"length(100);null"`
}
type Customer struct {
	*Model[Customer]
	CustomerId NumberField[uint64] `db:"primaryKey;autoIncrement"`
	Code       TextField           `db:"length(50);unique"` // Mã khách hàng
	Name       TextField           `db:"length(200)"`       // Tên khách
	Email      TextField           `db:"length(100);null"`
	Phone      TextField           `db:"length(20);null"`
	Address    TextField           `db:"length(300);null"`

	Note      TextField `db:"length(200);null"`
	CreatedAt DateTimeField
	UpdatedAt DateTimeField `db:"null"`
	CreatedBy TextField     `db:"length(100)"`
	UpdatedBy TextField     `db:"length(100);null"`
}
type Invoice struct {
	*Model[Invoice]
	InvoiceId       NumberField[uint64] `db:"primaryKey;autoIncrement"`
	OrderId         NumberField[uint64] `db:"index(order_invoice_idx)"` // Nếu liên kết tới Order
	Version         NumberField[int]    `db:"index(order_invoice_idx)"`
	CustomerId      NumberField[uint64] `db:"index(customer_invoice_idx)"`
	PaymentMethodId NumberField[uint64] `db:"index(payment_method_idx)"` // Mới thêm
	InvoiceDate     DateTimeField
	TotalAmount     NumberField[float64]

	Note      TextField `db:"length(200);null"`
	CreatedAt DateTimeField
	UpdatedAt DateTimeField `db:"null"`
	CreatedBy TextField     `db:"length(100)"`
	UpdatedBy TextField     `db:"length(100);null"`
}

type Item struct {
	*Model[Item]
	ItemId NumberField[uint64]  `db:"primaryKey;autoIncrement"`
	Code   TextField            `db:"length(50);unique"`
	Name   TextField            `db:"length(200)"`
	Unit   TextField            `db:"length(20)"` // Đơn vị tính
	Price  NumberField[float64] // Giá mặc định
	Note   TextField            `db:"length(200);null"`

	CreatedAt DateTimeField
	UpdatedAt DateTimeField `db:"null"`
	CreatedBy TextField     `db:"length(100)"`
	UpdatedBy TextField     `db:"length(100);null"`
}
type InvoiceDetail struct {
	*Model[InvoiceDetail]
	Id NumberField[uint64] `db:"primaryKey;autoIncrement"`

	InvoiceId NumberField[uint64] `db:"index(invoice_detail_idx)"`
	Version   NumberField[int]    `db:"index(invoice_detail_idx)"`

	ItemId NumberField[uint64]  `db:"index"`
	Amount NumberField[int]     // Số lượng
	Price  NumberField[float64] // Đơn giá
	Note   TextField            `db:"length(200);null"`

	CreatedAt DateTimeField
	UpdatedAt DateTimeField `db:"null"`
	CreatedBy TextField     `db:"length(100)"`
	UpdatedBy TextField     `db:"length(100);null"`
}
type OrderRepository struct {
	*Base
	Orders         *Order
	OrderItems     *OrderItem
	Invoices       *Invoice
	Customers      *Customer
	Items          *Item
	InvoiceDetails *InvoiceDetail
	PaymentMethods *PaymentMethod
}
type PaymentMethod struct {
	*Model[PaymentMethod]
	PaymentMethodId NumberField[uint64] `db:"primaryKey;autoIncrement"`
	Code            TextField           `db:"length(50);unique"` // Ví dụ: "CASH", "BANK", "CARD"
	Name            TextField           `db:"length(100)"`       // Ví dụ: "Tiền mặt", "Chuyển khoản"

	Note      TextField `db:"length(200);null"`
	CreatedAt DateTimeField
	UpdatedAt DateTimeField `db:"null"`
	CreatedBy TextField     `db:"length(100)"`
	UpdatedBy TextField     `db:"length(100);null"`
}

func TestIsSpecialChar(t *testing.T) {
	e := expression{}
	assert.True(t, e.isSpecialChar('('))
	assert.True(t, e.isSpecialChar(')'))
	assert.True(t, e.isSpecialChar(','))
	assert.True(t, e.isSpecialChar('.'))
	assert.True(t, e.isSpecialChar(' '))
	assert.True(t, e.isSpecialChar('/'))
	assert.True(t, e.isSpecialChar('+'))
	assert.True(t, e.isSpecialChar('-'))
	assert.True(t, e.isSpecialChar('*'))
	assert.True(t, e.isSpecialChar('%'))
	assert.True(t, e.isSpecialChar('='))
	assert.True(t, e.isSpecialChar('<'))
	assert.True(t, e.isSpecialChar('>'))
	assert.True(t, e.isSpecialChar('!'))
	assert.True(t, e.isSpecialChar('&'))
	assert.True(t, e.isSpecialChar('|'))
	assert.True(t, e.isSpecialChar('^'))
	assert.True(t, e.isSpecialChar('~'))
	assert.True(t, e.isSpecialChar('?'))
	assert.True(t, e.isSpecialChar(':'))
	assert.True(t, e.isSpecialChar(';'))
	assert.True(t, e.isSpecialChar('['))
	assert.True(t, e.isSpecialChar(']'))
	assert.True(t, e.isSpecialChar('{'))
	assert.True(t, e.isSpecialChar('}'))
	assert.True(t, e.isSpecialChar('@'))
	assert.True(t, e.isSpecialChar('#'))
	assert.True(t, e.isSpecialChar('$'))
}

func TestGetMarkList(t *testing.T) {
	dataTest := map[string][][]int{
		"select.(12/  select/12)": {[]int{0, 6}, {13, 19}},
		"select.(select)":         {[]int{0, 6}, {8, 14}},
		"max(select.select)":      {[]int{4, 10}, {11, 17}},
		"select.select":           {[]int{0, 6}, {7, 13}},
		"select.  select":         {[]int{0, 6}, {9, 15}},
	}
	e := expression{}
	for input, data := range dataTest {

		mark, err := e.getMarkList(input, "select")
		assert.NoError(t, err)
		assert.Equal(t, data, mark)

	}
}
func BenchmarkGetMarkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dataTest := map[string][][]int{
			"select.(12/  select/12)": {[]int{0, 6}, {13, 19}},
			"select.(select)":         {[]int{0, 6}, {8, 14}},
			"max(select.select)":      {[]int{4, 10}, {11, 17}},
			"select.select":           {[]int{0, 6}, {7, 13}},
			"select.  select":         {[]int{0, 6}, {9, 15}},
		}
		e := expression{}
		for input, data := range dataTest {

			mark, err := e.getMarkList(input, "select")
			assert.NoError(b, err)
			assert.Equal(b, data, mark)

		}
	}
}
func TestPrepare(t *testing.T) {
	dataTest := []string{
		"max(select.select)->max(`select`.`select`)",
		"max(select.select   + 10)->max(`select`.`select`   + 10)",
		"max(  select.  select   + 10)->max(  `select`.  `select`   + 10)",
		"select/select->`select`/`select`",
	}
	e := expression{}
	for _, data := range dataTest {
		input := strings.Split(data, "->")[0]
		expected := strings.Split(data, "->")[1]
		res, err := e.Prepare(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	}
}

func mssql() DialectCompiler {
	return &MssqlDialect
}
func TestCompile(t *testing.T) {
	funcTest(t)
}
func funcTest(t assert.TestingT) {
	Repository[OrderRepository]()
	e := expression{

		dialect: mssql(),
	}
	cmd := "ORDER.OrderID,Order.Note"

	compiled, err := e.compileSelect(cmd)

	assert.NoError(t, err)
	compiledExpected := "[orders].[order_id] AS [OrderID], [orders].[note] AS [Note]"
	assert.Equal(t, compiledExpected, compiled)

}
func TestCompileWithFuncCall(t *testing.T) {
	Repository[OrderRepository]()
	e := expression{

		dialect: mssql(),
	}
	cmd3 := "MAX(ORDER.OrderID, 10)"
	compiled3, err := e.compileSelect(cmd3)
	assert.NoError(t, err)
	compiledExpected3 := "MAX([orders].[order_id], 10)"
	assert.Equal(t, compiledExpected3, compiled3)
	cmd2 := "MAX(ORDER.OrderID)"
	compiled2, err := e.compileSelect(cmd2)
	assert.NoError(t, err)
	compiledExpected2 := "MAX([orders].[order_id])"
	assert.Equal(t, compiledExpected2, compiled2)

	cmd4 := "MAX(ORDER.OrderID, 10, 20)"
	compiled4, err := e.compileSelect(cmd4)
	assert.NoError(t, err)
	compiledExpected4 := "MAX([orders].[order_id], 10, 20)"
	assert.Equal(t, compiledExpected4, compiled4)
}

func BenchmarkCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcTest(b)
	}
}
