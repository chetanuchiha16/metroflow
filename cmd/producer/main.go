package main

import (
	"fmt"
	"metroflow/internal/producer"
	"net"
)


func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	producer.StreamLog(conn)
}