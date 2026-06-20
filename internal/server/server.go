package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

type server struct {
	ln net.Listener
}

func NewServer() *server {
	return &server{}
}

func (sv *server) StartServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	sv.ln = ln
	fmt.Printf("server started at %v\n", ln.Addr().String())
	sv.AcceptLoop()
}

func (sv *server) AcceptLoop() {
	for {
		conn, err := sv.ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("client connected at: %v\n", conn.RemoteAddr().String())
		buffer := fmt.Appendf(nil, "connected to %v\n", conn.LocalAddr().String())
		conn.Write(buffer)
		go sv.HandleConn(conn)
	}
}

func (sv *server) HandleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	sv.ReadLoop(conn, buffer)
}

func (sv *server) ReadLoop(conn net.Conn, buffer []byte) {
	jobs := make(chan string)
	workers := 1
	sv.fanOut(jobs, workers)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			// fmt.Printf("client %v disconnected %v\n", conn.RemoteAddr().String(), err)
			log.Fatalf("client %v disconnected %v\n", conn.RemoteAddr().String(), err)
			return
		}
		// fmt.Println(string(buffer[:n]))
		go func() {
			jobs <- string(buffer[:n])
		}()
	}
}

func (sv *server) fanOut(jobs <-chan string, workers int) {
	for worker := range workers {
		go func(worker int) {
			for job := range jobs {
				fmt.Printf("[%v] job \"%v\" by the worker %v\n", time.Now(), job, worker)
				time.Sleep(3 * time.Second)
			}
		}(worker)
	}
}
