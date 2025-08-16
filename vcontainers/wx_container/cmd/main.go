package main

import (
	"vdi"
	"wx_container"
)

// const configPath = "./cmd/config.yaml"

const configPath = "./../cmd/config.yaml"

func main() {

	vdi.Start[wx_container.WxContainer]()

}
