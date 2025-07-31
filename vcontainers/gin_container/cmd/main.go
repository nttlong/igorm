package main

import "gin_container"

const configPath = "./../cmd/config.yaml"

func main() {

	c := gin_container.NewGinContainer(configPath)
	c.StartServer()

}
