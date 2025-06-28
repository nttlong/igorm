package unvsef

import (
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
