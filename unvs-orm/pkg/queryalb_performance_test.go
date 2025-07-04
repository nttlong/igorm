package orm_test

import (
	"reflect"
	"testing"
	"time"
	orm "unvs-orm"
)

func BenchmarkQueryALB(b *testing.B) {
	var totalQueries int
	var totalTime int64

	for i := 0; i < b.N; i++ {
		start := time.Now()

		typ := reflect.TypeOf(&User{}).Elem()
		orm.Queryable[User](nil)
		tblName := orm.Utils.TableNameFromStruct(typ)
		retVal := orm.EntityUtils.QueryableFromType(typ, tblName, nil)
		retVal.Interface()

		elapsed := time.Since(start).Nanoseconds()
		totalTime += elapsed
		totalQueries++
	}

	b.ReportMetric(float64(totalQueries), "total_queries")
	b.ReportMetric(float64(totalTime)/float64(totalQueries), "avg_time_per_query_ns")
}
