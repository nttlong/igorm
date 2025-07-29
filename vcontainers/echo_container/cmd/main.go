package main

import "echo_container"

const configPath = "./../cmd/config.yaml"

func main() {

	c := echo_container.NewEchoContainer(configPath)
	c.StartServer()

}
