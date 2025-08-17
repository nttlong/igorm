package main

import (
	"log"
	apps "media/internal/services/app"
	"net/http"
	"os"
	"runtime/pprof"
	"wx"
)

func main() {

	err := wx.Start(func(app *apps.App) error {
		go func() {
			f, _ := os.Create("mem.pprof")
			pprof.WriteHeapProfile(f)
			f.Close()
			log.Println("pprof listening on :6060")
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
		server, err := app.Server.Ins()
		if err != nil {
			return err
		}
		server.Start()
		return nil
	})
	if err != nil {
		panic(err)
	}
}
