package main

import (
	"eks-inject/internal/server"
	"log"
)

func main() {
	log.Println("Starting server...")

	server.StartServer()
}
