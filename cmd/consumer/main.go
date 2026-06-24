package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	for {
		fmt.Println("waiting for job...")
		val, err := rdb.BLPop(context.Background(), 0, "job").Result()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(val)
	}
}
