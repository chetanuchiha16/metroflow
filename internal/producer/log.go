package producer

import (
	"fmt"
	"net"
	"time"
)

func StreamLog(conn net.Conn) {
	defer conn.Close()
	logStr := fmt.Sprintf("this is a log %v", 1)
	buffer := []byte(logStr)

	for {
		_, err := conn.Write(buffer)
		if err != nil {
			println(err)
			return
		}
		fmt.Println(string(logStr))
		time.Sleep(3 * time.Second)
	}
}
