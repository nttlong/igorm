package unvsef

import (
	"fmt"
	"testing"
)

type Comment struct {
	Id        Field[uint64] `db:"primaryKey;autoIncrement"`
	ArticleId Field[uint64] `db:"index"`
	Content   Field[string] `db:"FTS(content_idx)"`
}

type Repository struct {
	TenantDb
	// Articles *Article
	Comments *Comment
}
type Article struct {
	Id         FieldUint64    `db:"primaryKey;autoIncrement"`
	Title      FieldString    `db:"FTS(title_idx);length(50)"`
	Content    FieldString    `db:"FTS(content_idx);length(50)"`
	CreatedOn  FieldDateTime  `db:"default:now()"`
	CreatedBy  FieldString    `db:"length(50)"`
	ModifiedOn *FieldDateTime `db:"default:now();"`
	ModifiedBy *FieldString   `db:"length(50)"`
}

func TestRepository(t *testing.T) {

	qr := Queryable[Article]()
	d := &PostgresDialect{}
	sql, args := "", []interface{}{}
	wsum := qr.Content.Len().Sum()
	sql, args = wsum.ToSqlExpr(d)
	t.Log(sql, args)
	wsum = qr.Id.Sum()
	sql, args = wsum.ToSqlExpr(d)
	t.Log(sql, args)
	w1 := qr.Content.Len()
	sql, args = w1.ToSqlExpr(d)
	t.Log(sql, args)
	wid := qr.Id.Eq(1)
	sql, args = wid.ToSqlExpr(d)
	t.Log(sql, args)
	//sql, args := qr.Content.Len().Eq(qr.CreatedBy).ToSqlExpr(d)

	where := qr.Content.Len().Eq(qr.CreatedBy)
	sql, args = where.ToSqlExpr(d)
	//sql= "[articles].[content] = ?"
	//mong muon la "[articles].[content] = [articles].[CreatedBy]"
	t.Log(sql, args)
	bw1 := qr.CreatedOn.Between(2020, 2025)
	sql, args = bw1.ToSqlExpr(d)

	//sql= "[articles].[content] = ?"
	//mong muon la "[articles].[content] = [articles].[CreatedBy]"
	t.Log(sql, args)
	bw2 := (qr.ModifiedOn.Between(2020, 2025))
	sql, args = bw1.ToSqlExpr(d)

	//sql= "[articles].[content] = ?"
	//mong muon la "[articles].[content] = [articles].[CreatedBy]"
	t.Log(sql, args)
	bw := bw1.And(bw2)
	sql, args = bw.ToSqlExpr(d)
	t.Log(sql, args)
}
func TestQuery(t *testing.T) {
	d := &PostgresDialect{}
	article := Queryable[Article]()
	comment := Queryable[Comment]()
	joinExpr, args := compiler.Compile(comment.ArticleId.Eq(article.Id), d)
	fmt.Println(joinExpr, args)
	sql := From(article).Select(article.Content, article.CreatedBy).Where(article.Content.Len().Gt(100))
	sqlStr, args := sql.ToSQL(d)
	t.Log(sqlStr, args)
	conditional := comment.ArticleId.Eq(article.Id).And(comment.Content.Len().Gt(100))
	joinExpr, args = compiler.Compile(conditional, d)
	fmt.Println(joinExpr, args)
	//jonnExpr la
	//"((\"comments\".\"article_id\" = ?) AND (LEN(\"comments\".\"content\") > ?))"
	// va args la 100

	sql3 := From(LeftJoin(conditional)).Select(comment.Content)
	sqlStr, args = sql3.ToSQL(d)
	t.Log(sqlStr, args)

}
