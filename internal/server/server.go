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
		fmt.Println("waiting for handleconn to finish")
		wg.Wait()
		fmt.Println("waiting for handle conn finished")
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
	// wg.Add(1)
	sv.fanOut(jobs, workers)
	// var wg sync.WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			line, _, err := reader.ReadLine() // back pressure, i think
			if err != nil {
				log.Printf("client %v disconnected %v\n", conn.RemoteAddr().String(), err) // so program exit after the client leaves, to track time
				// sv.exit <- true
				return
			}
			// wg.Go(func() {
			jobs <- string(line)
			// })
			// wg.Wait()
		}
	}(&wg)
	fmt.Println("waiting for sending jobs...")
	wg.Wait() 
	fmt.Println("waiting for sending jobs finished")
	sv.exit <- true

}

func (sv *server) fanOut(jobs <-chan string, workers int) {
	var wg sync.WaitGroup
	for worker := range workers {
		wg.Add(1)
		go func(worker int, wg *sync.WaitGroup) {
			defer wg.Done()
			for job := range jobs {
				// start := time.Now()
				fmt.Printf("[%v] job \"%v\" being done by the worker %v...\n", time.Now().Format(time.TimeOnly), job, worker)
				time.Sleep(3 * time.Second) // simulate processing
				fmt.Printf("worker %v finished the job.\n", worker)
			}
		}(worker, &wg)
	}
	fmt.Println("waiting for fanout to finish...")
	wg.Wait() // waiting for fanout to finish and fanout() is syncronous so this is blocking, job is never being sent deadlock
	fmt.Println("waiting for fanout finished")
}
