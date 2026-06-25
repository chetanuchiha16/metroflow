package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"metroflow/pkg"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	fmt.Println("waiting for job...")
	mp := make(map[string]int)
	file, err := os.Create("log_count.csv")
	if err != nil {
		fmt.Printf("error opening the csv file error: %v\n", err)
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"Time stamp", "ERROR", "WARN", "INFO"})
	
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
		if mp["ERROR"] > 1 {
			go func() {
				defer csvWriter.Flush()
				row := []string{time.Now().Format(time.TimeOnly), fmt.Sprint(mp["ERROR"]), fmt.Sprint(mp["WARN"]), fmt.Sprint(mp["INFO"])}
				csvWriter.Write(row)

			}()
		}
		fmt.Println(mp)
	}
}
