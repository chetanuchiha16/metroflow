package main

import (
	"metroflow/internal/server"
)

func main() {
	server := server.NewServer()
	server.StartServer()
}


