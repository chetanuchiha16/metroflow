package main

import (
	"log"
	"metroflow/internal/producer"
	"net"
)


func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	producer.StreamLog(conn)
}