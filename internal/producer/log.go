package producer

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func StreamLog(conn net.Conn) {
	defer conn.Close()
	// i := 0
	levels := []string {"INFO", "WARN", "ERROR"}
	for i := range 10 {

		logStr := fmt.Sprintf("this is a %v log %v\n",levels[rand.Intn(i + 1) % len(levels)], i)
		buffer := []byte(logStr)
		n, err := conn.Write(buffer)
		if err != nil {
			println(err)
			return
		}
		fmt.Println(string(logStr))
		fmt.Println(n)
		time.Sleep(1 * time.Millisecond)
		// i++
	}
}
