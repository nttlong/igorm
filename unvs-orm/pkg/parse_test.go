package orm_test

import (
	"testing"
	EXPR "unvs-orm/expr"

	"github.com/stretchr/testify/assert"
)

// QuoteSqlIdentifiers adds backticks to table names, column names, and aliases in a SQL string.

func BenchmarkTestParse(t *testing.B) {
	cfg := EXPR.NewExprConfig()
	for i := 0; i < t.N; i++ {
		sql11 := "items.basicSalary*100+min(bonus, 1000)*[as]/len('len(as)'+'group' as name) AS total"
		sql11 = cfg.QuoteExpression(sql11)
		expected11 := "`items`.`basicSalary`*100+min(`bonus`, 1000)*`as`/len('len(as)'+'group' as `name`) AS `total`"
		assert.Equal(t, expected11, sql11)

		sql7 := "items.basicSalary*100+min(bonus, 1000)"
		sql7 = cfg.QuoteExpression(sql7)
		expected7 := "`items`.`basicSalary`*100+min(`bonus`, 1000)"
		assert.Equal(t, expected7, sql7)

		sql7 = "items.basicSalary*100+min(bonus, 1000)"
		sql7 = cfg.QuoteExpression(sql7)
		expected7 = "`items`.`basicSalary`*100+min(`bonus`, 1000)"
		assert.Equal(t, expected7, sql7)

		sql8 := "items.basicSalary*100+min(bonus, 1000)/len(items.name)"
		sql8 = cfg.QuoteExpression(sql8)
		expected8 := "`items`.`basicSalary`*100+min(`bonus`, 1000)/len(`items`.`name`)"
		assert.Equal(t, expected8, sql8)
		sql9 := "items.basicSalary*100+min(bonus, 1000)/len(items.name) as total"
		sql9 = cfg.QuoteExpression(sql9)
		expected9 := "`items`.`basicSalary`*100+min(`bonus`, 1000)/len(`items`.`name`) as `total`"
		assert.Equal(t, expected9, sql9)
		sql10 := "items.basicSalary*100+min(bonus, 1000)/len(items.as) as total"
		sql10 = cfg.QuoteExpression(sql10)
		expected10 := "`items`.`basicSalary`*100+min(`bonus`, 1000)/len(`items`.`as`) as `total`"
		assert.Equal(t, expected10, sql10)

	}
}
func TestParse(t *testing.T) {
	cfg := EXPR.NewExprConfig()
	sql7 := "items.basicSalary*100+min(bonus, 1000)"
	sql7 = cfg.QuoteExpression(sql7)
	expected7 := "`items`.`basicSalary`*100+min(`bonus`, 1000)"
	assert.Equal(t, expected7, sql7)

	sql7 = "items.basicSalary*100+min(bonus, 1000)"
	sql7 = cfg.QuoteExpression(sql7)
	expected7 = "`items`.`basicSalary`*100+min(`bonus`, 1000)"
	assert.Equal(t, expected7, sql7)

	sql8 := "items.basicSalary*100+min(bonus, 1000)/len(items.name)"
	sql8 = cfg.QuoteExpression(sql8)
	expected8 := "`items`.`basicSalary`*100+min(`bonus`, 1000)/len(`items`.`name`)"
	assert.Equal(t, expected8, sql8)
	sql9 := "items.basicSalary*100+min(bonus, 1000)/len(items.name) as total"
	sql9 = cfg.QuoteExpression(sql9)
	expected9 := "`items`.`basicSalary`*100+min(`bonus`, 1000)/len(`items`.`name`) as `total`"
	assert.Equal(t, expected9, sql9)
	sql10 := "items.basicSalary*100+min(bonus, 1000)/len(items.as) as total"
	sql10 = cfg.QuoteExpression(sql10)
	expected10 := "`items`.`basicSalary`*100+min(`bonus`, 1000)/len(`items`.`as`) as `total`"
	assert.Equal(t, expected10, sql10)
	sql11 := "items.basicSalary*100+min(bonus, 1000)*[as]/len('len(as)'+'group' as name) AS total"
	sql11 = cfg.QuoteExpression(sql11)
	expected11 := "`items`.`basicSalary`*100+min(`bonus`, 1000)*`as`/len('len(as)'+'group' as `name`) AS `total`"
	assert.Equal(t, expected11, sql11)

}
