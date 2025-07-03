package orm_test

import (
	"reflect"
	"testing"
	orm "unvs-orm"
)

func BenchmarkQueryALB(b *testing.B) {
	var totalQueries int
	var totalTime int64

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		typ := reflect.TypeOf(&User{}).Elem()
		orm.Queryable[User]()
		tblName := orm.Utils.TableNameFromStruct(typ)
		retVal := orm.EntityUtils.QueryableFromType(typ, tblName)
		retVal.Interface()
		b.StopTimer()

		totalQueries++
		totalTime += b.Elapsed().Nanoseconds()
	}

	b.ReportMetric(float64(totalQueries), "total_queries")
	b.ReportMetric(float64(totalTime)/float64(totalQueries), "avg_time_per_query_ns")
}
