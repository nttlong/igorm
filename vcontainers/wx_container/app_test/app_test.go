package app_test

import (
	"testing"
	"vdi"
	"wx_container"
)

func TestAppStart(t *testing.T) {
	 vdi.Start[wx_container.WxContainer]()
}
