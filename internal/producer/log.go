package producer

import (
	"fmt"
	"net"
	"time"
)

func StreamLog(conn net.Conn) {
	defer conn.Close()
	// i := 0

	for i := range 10 {
		logStr := fmt.Sprintf("this is a log %v", i)
		buffer := []byte(logStr)
		_, err := conn.Write(buffer)
		if err != nil {
			println(err)
			return
		}
		fmt.Println(string(logStr))
		time.Sleep(1 * time.Millisecond)
		// i++
	}
}
