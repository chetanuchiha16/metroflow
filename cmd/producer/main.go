package main

import (
	"fmt"
	"log"
	"metroflow/internal/producer"
	"net"
)


func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("connected to %v\n", conn.RemoteAddr().String())
	producer.StreamLog(conn)
}