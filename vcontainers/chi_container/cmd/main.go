package main

import "chi_container"

const configPath = "./../cmd/config.yaml"

func main() {

	c := chi_container.NewChiContainer(configPath)
	c.StartServer()

}
