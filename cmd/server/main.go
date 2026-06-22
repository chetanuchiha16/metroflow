package main

import (
	"fmt"
	"metroflow/internal/server"
	// "os"

	"os"
	"time"
)

func main() {
	start := time.Now()
	exit := make(chan bool)
	go func() { // wait for exit before starting, else this go routine hasnt even started
		if (<-exit) {
			// os.Exit(0)
			fmt.Printf("all jobs finished in %v\n", time.Since(start)) // does this really mean all jobs ? or the handle connection one ?
			os.Exit(0)
		}
	}()
	server := server.NewServer(exit)
	server.StartServer()
	

}


