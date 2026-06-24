package main

import (
	"fmt"
	"metroflow/internal/server"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	start := time.Now()
	server := server.NewServer(rdb)
	server.StartServer()
	fmt.Printf("all jobs finished in %v\n", time.Since(start)) // does this really mean all jobs ? or the handle connection one ?

}
