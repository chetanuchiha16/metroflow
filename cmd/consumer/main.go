package main

import (
	"context"
	"encoding/json"
	"fmt"
	"metroflow/pkg"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	fmt.Println("waiting for job...")
	mp := make(map[string]int)
	for {
		job, err := rdb.BLPop(context.Background(), 0, "job").Result()
		if err != nil {
			fmt.Println(err)
			break
		}
		valBytes := []byte(job[1])
		var jsonLog pkg.LogType
		if err = json.Unmarshal(valBytes, &jsonLog); err != nil {
			fmt.Printf("error unmarshel %v", err)
		}
		mp[jsonLog.Level] += 1
		fmt.Println(mp)
	}
}
