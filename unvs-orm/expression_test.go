package orm

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlparser "github.com/xwb1989/sqlparser"
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

type expression struct {
	dialect     DialectCompiler
	cmp         *CompilerUtils
	keywords    []string
	specialChar []byte
}

func (e *expression) getField(expr interface{}) (string, error) {

	switch expr := expr.(type) {
	case *sqlparser.ColName:
		tableNameFromSyntax := expr.Qualifier.ToViewName().Name.CompliantName()
		tableName := tableNameFromSyntax
		if dbTableName := Utils.GetDbTableName(tableNameFromSyntax); dbTableName != "" {
			tableName = dbTableName
		}
		if tableName == "" {
			return "", fmt.Errorf("table name not found for %s", expr)
		}
		fieldName := expr.Name.CompliantName()
		metaInfo := Utils.GetMetaInfoByTableName(tableName)
		if metaInfo != nil {
			if _, ok := metaInfo[strings.ToLower(fieldName)]; !ok {
				return "", fmt.Errorf("field %s not found in table %s", fieldName, tableName)
			} else {
				fieldName = Utils.ToSnakeCase(fieldName)
			}
		}
		ret := e.cmp.Quote(tableName, fieldName) + " AS " + e.cmp.Quote(expr.Name.CompliantName())
		return ret, nil

	case *sqlparser.AliasedExpr:
		return e.getField(expr.Expr)
	default:
		return "", fmt.Errorf("not support %s", expr)
	}
	return "", nil

}
func (e *expression) InsertMarks(input string, markList [][]int) string {
	// Sort theo vị trí start giảm dần để tránh làm lệch chỉ số khi chèn
	sort.Slice(markList, func(i, j int) bool {
		return markList[i][0] > markList[j][0]
	})

	builder := strings.Builder{}
	builder.WriteString(input)

	for _, mark := range markList {
		start, end := mark[0], mark[1]
		if start < 0 || end > builder.Len() || start >= end {
			continue // tránh lỗi khi input sai
		}

		// Chèn dấu '`' vào cuối đoạn trước
		builderStr := builder.String()
		builder.Reset()
		builder.WriteString(builderStr[:end])
		builder.WriteString("`")
		builder.WriteString(builderStr[end:])

		builderStr = builder.String()
		builder.Reset()
		builder.WriteString(builderStr[:start])
		builder.WriteString("`")
		builder.WriteString(builderStr[start:])
	}

	return builder.String()
}
func (e *expression) isSpecialChar(input byte) bool {
	if e.specialChar == nil {
		e.specialChar = []byte{
			'(', ')', ',', '.', ' ', '/', '+', '-', '*', '%', '=', '<', '>', '!', '&', '|', '^', '~', '?', ':', ';', '[', ']', '{', '}', '@', '#', '$',
		}
	}
	return bytes.Contains(e.specialChar, []byte{input})
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
func (e *expression) getMarkList(input string, keyword string) ([][]int, error) {
	keyword = strings.ToLower(keyword)
	if !strings.Contains(strings.ToLower(input), keyword) {
		return nil, nil
	}
	beginKeyword := keyword[0]
	beginKeywordUpper := strings.ToUpper(string(beginKeyword))[0]
	input = strings.TrimLeft(input, " ")
	input = strings.TrimRight(input, " ")

	check := ""

	start := 0
	end := -1
	markList := [][]int{}
	if !strings.Contains(strings.ToLower(input), keyword) {
		return markList, nil
	}
	i := 0

	for i < len(input) {

		check += strings.ToLower(string(input[i]))
		if e.isSpecialChar(input[i]) {
			if i+1 > len(input) {

				return nil, fmt.Errorf("syntax error in \n%s", input)
			} else {
				for i+1 < len(input) && (input[i+1] != beginKeyword && input[i+1] != beginKeywordUpper) {
					i++

				}
				if i+1 >= len(input) {
					return markList, nil
				}
				check = strings.ToLower(string(input[i+1]))

				start = i + 1
				i++
			}
		}
		if check == keyword {
			j := i

			for j+1 < len(input) && !(e.isSpecialChar(input[j+1])) {

				j++
			}
			if j < len(input) {
				end = i

				markList = append(markList, []int{start, end + 1})
				check = ""
				start = end + 1
				end = -1
				i = j
			} else {
				i++
			}
		}
		i++
	}
	return markList, nil
}
func (e *expression) Prepare(input string) (string, error) {
	if e.keywords == nil {
		e.keywords = []string{
			"select",
			"from",
			"where",
			"group",
			"order",
			"limit",
			"offset",
		}
	}
	for _, keyword := range e.keywords {
		markList, err := e.getMarkList(input, keyword)
		if err != nil {
			return "", err
		}
		input = e.InsertMarks(input, markList)
	}
	return input, nil

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

func (e *expression) compileSelect(cmd string) (string, error) {
	if e.cmp == nil {
		e.cmp = e.dialect.getCompiler()
	}
	cmd, err := e.Prepare(cmd)
	if err != nil {
		return "", err
	}
	sqlTest := "select " + cmd
	stm, err := sqlparser.Parse(sqlTest)
	if err != nil {
		return "", err
	}
	fields := []string{}
	if stmt, ok := stm.(*sqlparser.Select); ok {
		for _, col := range stmt.SelectExprs {
			fieldE, err := e.getField(col)
			if err != nil {
				return "", err
			}
			fields = append(fields, fieldE)
		}
	} else {
		return "", fmt.Errorf("%s not a select statement", cmd)
	}

	return strings.Join(fields, ", "), nil
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

	cmd2 := "select.select"
	compiled2, err := e.compileSelect(cmd2)
	assert.NoError(t, err)
	compiledExpected2 := "MAX([orders].[order_id])"
	assert.Equal(t, compiledExpected2, compiled2)
}

func BenchmarkCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcTest(b)
	}
}
