package main

import (
	"dbv"
	"os"
	"runtime/pprof"
)

func main() {
	//go tool pprof -http=localhost:8080  ./cpu1-eorm.prof
	f, _ := os.Create("cpu1-eorm.prof")
	//go tool pprof -http=localhost:8080 ./cpu19-ent.prof
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < 1000000; i++ {
		joinExpr := "Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'"
		builder := dbv.SqlBuilder.From(joinExpr).Select()
		builder.ToSql(dbv.DialectFactory.Create("mssql"))
	}
}
