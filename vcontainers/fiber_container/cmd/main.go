package main

import "fiber_container"

const configPath = "./../cmd/config.yaml"

func main() {

	c := fiber_container.NewFiberContainer(configPath)
	c.StartServer()

}
