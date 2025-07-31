package main

import "http_container"

const configPath = "./../cmd/config.yaml"

func main() {

	c := http_container.NewHttpContainer(configPath)
	c.StartServer()

}
