package main

import (
	"context"
	"fmt"
	"metroflow/internal/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	start := time.Now()
	server := server.NewServer(rdb, ctx)
	server.StartServer()
	fmt.Printf("all jobs finished in %v\n", time.Since(start)) // does this really mean all jobs ? or the handle connection one ?

}
