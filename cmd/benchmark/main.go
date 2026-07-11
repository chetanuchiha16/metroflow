package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Read server greeting
	reader := bufio.NewReader(conn)
	greeting, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server greeting: %s", greeting)

	const totalLogs = 100000
	fmt.Printf("Streaming %d logs...\n", totalLogs)

	start := time.Now()
	writer := bufio.NewWriter(conn)
	for i := 0; i < totalLogs; i++ {
		_, err := writer.WriteString("this is a INFO log 0\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	duration := time.Since(start)
	fmt.Printf("Sent %d logs in %v (Throughput: %.2f logs/sec)\n", 
		totalLogs, duration, float64(totalLogs)/duration.Seconds())
}
