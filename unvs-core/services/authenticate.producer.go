package services

import (
	"dbx"
	"fmt"
	"time"

	"unvs.core/models"
)

var loginChange = make(chan LoginChanItem, 1000)

type LoginChanItem struct {
	info *models.LoginInfo
	db   *dbx.DBXTenant
}

func (u *AuthenticateService) producer(info *models.LoginInfo) {
	item := LoginChanItem{
		info: info,
		db:   u.TenantDb.Clone("loginInfoDB"),
	}
	select {
	case loginChange <- item:
	default:
		fmt.Println("loginChange channel is full, dropping item")
	}
}
func consumer() {
	for {
		item := <-loginChange
		dataInsert := item.info
		err := item.db.Insert(*dataInsert)
		fmt.Println("inserted login info", *dataInsert)

		if err != nil {
			fmt.Println("error inserting login info", err)
		}
		_, err = item.db.Update(&models.User{}).Where(
			"UserId = ?", item.info.UserId).Set(
			"LastLoginAt", time.Now(),
		).Execute()
		if err != nil {
			fmt.Println("error updating user last login at", err)
		} else {

			fmt.Println("updated user last login")
		}

	}
}
func init() {
	go consumer() // run consumer in background
}
