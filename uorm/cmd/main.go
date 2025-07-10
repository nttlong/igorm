package main

import (
	"os"
	"runtime/pprof"
	"uorm"
)

type User struct {
	uorm.Model
	UserId   uorm.Field
	UserName uorm.Field
	Email    uorm.Field
}

func main() {
	// go tool pprof -http=localhost:8080  cpu11.prof
	f, _ := os.Create("cpu8.prof")
	//go tool pprof -http=localhost:8080 ./cpu19-ent.prof
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < 1000000; i++ {
		qr := uorm.Queryable[User](uorm.DB_TYPE_MSSQL, "users")

		selector := qr.Selector(qr.Email.Add("qr.UserId"), qr.UserName.Add("qr.UserName"))
		selector.String()
	}
}
