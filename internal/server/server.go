package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	// "sync"
	"time"
)

type server struct {
	ln   net.Listener
	exit chan bool
}

func NewServer(exit chan bool) *server {
	return &server{exit: exit}
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
		var wg sync.WaitGroup
		start := time.Now()
		wg.Add(1)
		go sv.HandleConn(conn, &wg)
		wg.Wait()
		fmt.Printf("handle connection finished in %v\n", time.Since(start))
	}
}

func (sv *server) HandleConn(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()
	sv.ReadLoop(conn)
}

func (sv *server) ReadLoop(conn net.Conn) {
	reader := bufio.NewReader(conn)
	jobs := make(chan string)
	workers := 10
	sv.fanOut(jobs, workers)
	// var wg sync.WaitGroup
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Printf("client %v disconnected %v\n", conn.RemoteAddr().String(), err) // so program exit after the client leaves, to track time
			sv.exit <- true
			return
		}
		// wg.Go(func() {
		jobs <- string(line)
		// })
		// wg.Wait()
	}
}

func (sv *server) fanOut(jobs <-chan string, workers int) {
	for worker := range workers {
		go func(worker int) {
			for job := range jobs {
				// start := time.Now()
				fmt.Printf("[%v] job \"%v\" by the worker %v\n", time.Now(), job, worker)
				time.Sleep(3 * time.Second) // simulate processing
			}
		}(worker)
	}
}
