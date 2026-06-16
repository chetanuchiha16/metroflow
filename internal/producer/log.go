package producer

import (
	"fmt"
	"net"
	"time"
)

func StreamLog(conn net.Conn) string {
	defer conn.Close()
	log := fmt.Sprintf("this is a log %v", 1)
	buffer := []byte(log)

	for {
		n, err := conn.Write(buffer)
		if err != nil {
			println(err)
		}
		fmt.Println(string(buffer[:n]))
		time.Sleep(3 * time.Second)
	}
}
