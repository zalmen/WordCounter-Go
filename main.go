package main

import (
	"yuval/home-exercise/cyolo/httpserver"
)

func main() {
	server := httpserver.New()
	server.Start()
}
