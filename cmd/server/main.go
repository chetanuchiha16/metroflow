package main

import (
	"fmt"
	"metroflow/internal/server"
	"time"
)

func main() {
	start := time.Now()
	server := server.NewServer()
	server.StartServer()
	fmt.Printf("all jobs finished in %v\n", time.Since(start)) // does this really mean all jobs ? or the handle connection one ?
	

}


