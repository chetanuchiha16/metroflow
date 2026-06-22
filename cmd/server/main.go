package main

import (
	// "fmt"
	"metroflow/internal/server"
	// "os"

	// "os"
	// "time"
)

func main() {
	
	// fmt.Printf("all jobs finished in %v\n", time.Since(start)) // does this really mean all jobs ? or the handle connection one ?
	server := server.NewServer()
	server.StartServer()
	

}


