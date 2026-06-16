package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("server started at %v\n", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("client connected at: %v\n", conn.RemoteAddr().String())
		buffer := fmt.Appendf(nil, "connected to %v\n", conn.LocalAddr().String())
		conn.Write(buffer)
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("client %v disconnected %v\n", conn.RemoteAddr().String(), err)
			return
		}
		fmt.Println(string(buffer[:n]))
	}
}
