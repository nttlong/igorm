package orm

import (
	"strings"
	"testing"
	EXPR "unvs-orm/expr"

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
	e := EXPR.ExpressionTest{}
	assert.True(t, e.IsSpecialChar('('))
	assert.True(t, e.IsSpecialChar(')'))
	assert.True(t, e.IsSpecialChar(','))
	assert.True(t, e.IsSpecialChar('.'))
	assert.True(t, e.IsSpecialChar(' '))
	assert.True(t, e.IsSpecialChar('/'))
	assert.True(t, e.IsSpecialChar('+'))
	assert.True(t, e.IsSpecialChar('-'))
	assert.True(t, e.IsSpecialChar('*'))
	assert.True(t, e.IsSpecialChar('%'))
	assert.True(t, e.IsSpecialChar('='))
	assert.True(t, e.IsSpecialChar('<'))
	assert.True(t, e.IsSpecialChar('>'))
	assert.True(t, e.IsSpecialChar('!'))
	assert.True(t, e.IsSpecialChar('&'))
	assert.True(t, e.IsSpecialChar('|'))
	assert.True(t, e.IsSpecialChar('^'))
	assert.True(t, e.IsSpecialChar('~'))
	assert.True(t, e.IsSpecialChar('?'))
	assert.True(t, e.IsSpecialChar(':'))
	assert.True(t, e.IsSpecialChar(';'))
	assert.True(t, e.IsSpecialChar('['))
	assert.True(t, e.IsSpecialChar(']'))
	assert.True(t, e.IsSpecialChar('{'))
	assert.True(t, e.IsSpecialChar('}'))
	assert.True(t, e.IsSpecialChar('@'))
	assert.True(t, e.IsSpecialChar('#'))
	assert.True(t, e.IsSpecialChar('$'))
}

func TestGetMarkList(t *testing.T) {
	dataTest := map[string][][]int{
		"select.(12/  select/12)": {[]int{0, 6}, {13, 19}},
		"select.(select)":         {[]int{0, 6}, {8, 14}},
		"max(select.select)":      {[]int{4, 10}, {11, 17}},
		"select.select":           {[]int{0, 6}, {7, 13}},
		"select.  select":         {[]int{0, 6}, {9, 15}},
	}
	e := EXPR.ExpressionTest{}
	for input, data := range dataTest {

		mark, err := e.GetMarkList(input, "select")
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
		e := EXPR.ExpressionTest{}
		for input, data := range dataTest {

			mark, err := e.GetMarkList(input, "select")
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
	e := EXPR.ExpressionTest{}
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
func TestFieldExpr(t *testing.T) {
	e := EXPR.ExpressionTest{
		DbDriver: EXPR.DB_TYPE_MSSQL,
	}
	cmd := "ORDER.OrderID"
	tables := []string{}
	context := map[string]string{}
	compiled, err := e.Compile(&tables, &context, cmd, true)
	assert.NoError(t, err)
	compiledExpected := "[orders].[order_id]"
	assert.Equal(t, compiledExpected, compiled.Syntax)
	assert.Equal(t, []string{"orders"}, tables)
	assert.Equal(t, map[string]string{"orders": "T1"}, context)

}
func TestSumExpr(t *testing.T) {
	e := EXPR.ExpressionTest{
		DbDriver: EXPR.DB_TYPE_MSSQL,
	}
	cmd := "SUM(ORDER.OrderID)"
	tables := []string{}
	context := map[string]string{}
	compiled, err := e.Compile(&tables, &context, cmd, true)
	assert.NoError(t, err)
	compiledExpected := "SUM([orders].[order_id])"
	assert.Equal(t, compiledExpected, compiled.Syntax)
	assert.Equal(t, []string{"orders"}, tables)
	assert.Equal(t, map[string]string{"orders": "T1"}, context)

}
func funcTest(t assert.TestingT) {
	Repository[OrderRepository]()
	e := EXPR.ExpressionTest{
		DbDriver: EXPR.DB_TYPE_MSSQL,
	}
	cmd := "ORDER.OrderID,Order.Note"

	tables := []string{}
	context := map[string]string{}
	compiled, err := e.Compile(&tables, &context, cmd, true)

	assert.NoError(t, err)
	compiledExpected := "[orders].[order_id] AS [OrderID], [orders].[note] AS [Note]"
	assert.Equal(t, compiledExpected, compiled)

}
func TestCompileWithFuncCall(t *testing.T) {
	Repository[OrderRepository]()
	e := EXPR.ExpressionTest{
		DbDriver: EXPR.DB_TYPE_MSSQL,
	}
	tables := []string{}
	context := map[string]string{}

	cmd6 := "text(ORDER.OrderID)"
	compiled6, err := e.Compile(&tables, &context, cmd6, false)
	assert.NoError(t, err)
	compiledExpected6 := "CONVERT(NVARCHAR(50), [orders].[order_id])"
	assert.Equal(t, compiledExpected6, compiled6)

	cmd5 := "order.OrderID+order.Amount as TotalAmount"
	compiled5, err := e.Compile(&tables, &context, cmd5, false)
	assert.NoError(t, err)
	compiledExpected5 := "[orders].[order_id] + [orders].[order_id] AS [TotalAmount]"

	assert.Equal(t, compiledExpected5, compiled5)
	cmd3 := "MAX(ORDER.OrderID, 10)"
	compiled3, err := e.Compile(&tables, &context, cmd3, false)
	assert.NoError(t, err)
	compiledExpected3 := "MAX([orders].[order_id], 10)"
	assert.Equal(t, compiledExpected3, compiled3)
	cmd2 := "MAX(ORDER.OrderID)"
	compiled2, err := e.Compile(&tables, &context, cmd2, false)
	assert.NoError(t, err)
	compiledExpected2 := "MAX([orders].[order_id])"
	assert.Equal(t, compiledExpected2, compiled2)

	cmd4 := "MAX(ORDER.OrderID, 10, 20)"
	compiled4, err := e.Compile(&tables, &context, cmd4, false)
	assert.NoError(t, err)
	compiledExpected4 := "MAX([orders].[order_id], 10, 20)"
	assert.Equal(t, compiledExpected4, compiled4)

}

func BenchmarkCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcTest(b)
	}
}
