package main

import (
	"eks-injector/internal/server"
	"log"
)

func main() {
	log.Println("Starting server...")

	server.StartServer()
}
