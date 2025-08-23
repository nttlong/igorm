package main

import (
	"os"
	"runtime/pprof"
	"vdb"
)

func main() {
	//go tool pprof -http=localhost:8080  ./cpu1-eorm.prof
	f, _ := os.Create("cpu1-eorm.prof")
	//go tool pprof -http=localhost:8080 ./cpu19-ent.prof
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()
	for i := 0; i < 1000000; i++ {
		joinExpr := "Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'"
		builder := vdb.SqlBuilder.From(joinExpr).Select()
		builder.ToSql(vdb.DialectFactory.Create("mssql"))
	}
}
