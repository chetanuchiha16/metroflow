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
		fmt.Printf("[%v] for handleconn to finish\n", time.Now().Format(time.TimeOnly))
		wg.Wait()
		fmt.Printf("[%v] for handle conn finished\n", time.Now().Format(time.TimeOnly))
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
	var fanoutWg sync.WaitGroup
	fanoutWg.Add(1)
	go sv.fanOut(jobs, workers, &fanoutWg)
	// wg.Add(1)
	// var wg sync.WaitGroup
	var sendJobWg sync.WaitGroup
	sendJobWg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer close(jobs)
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
	}(&sendJobWg)
	fmt.Printf("[%v] for sending jobs...\n", time.Now().Format(time.TimeOnly))
	sendJobWg.Wait()
	fmt.Printf("[%v] for sending jobs finished.\n", time.Now().Format(time.TimeOnly))

	fmt.Printf("[%v] for fanout to finish...\n", time.Now().Format(time.TimeOnly))
	fanoutWg.Wait()
	fmt.Printf("[%v] for fanout finished.\n", time.Now().Format(time.TimeOnly))

	sv.exit <- true

}

func (sv *server) fanOut(jobs <-chan string, workers int, wg *sync.WaitGroup) {
	defer wg.Done()
	var workerWg sync.WaitGroup
	for worker := range workers {
		workerWg.Add(1)
		go func(worker int, workerWg *sync.WaitGroup) {
			defer workerWg.Done()
			for job := range jobs {
				// start := time.Now()
				fmt.Printf("[%v] job \"%v\" being done by the worker %v...\n", time.Now().Format(time.TimeOnly), job, worker)
				time.Sleep(3 * time.Second) // simulate processing
				fmt.Printf("[%v]worker %v finished the job %v.\n", time.Now().Format(time.TimeOnly), worker, job)
			}
		}(worker, &workerWg)
	}
	fmt.Printf("[%v] for all workers to finish...\n", time.Now().Format(time.TimeOnly))
	workerWg.Wait() // waiting for fanout to finish and fanout() is syncronous so this is blocking, job is never being sent deadlock ?
	fmt.Printf("[%v] for all workers finished\n", time.Now().Format(time.TimeOnly))
}
