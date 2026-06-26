package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"metroflow/pkg"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type server struct {
	ln  net.Listener
	rdb *redis.Client
	ctx context.Context
	wg  *sync.WaitGroup
}

func NewServer(ctx context.Context, wg *sync.WaitGroup, rdb *redis.Client) *server {
	return &server{
		ctx: ctx,
		wg:  wg,
		rdb: rdb,
	}
}

func (sv *server) StartServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	sv.ln = ln
	fmt.Printf("server started at %v\n", ln.Addr().String())
	sv.wg.Add(1)
	go func() {
		defer sv.wg.Done()
		<-sv.ctx.Done()
		fmt.Println("\nshutting down...")
		ln.Close()
	}()
	sv.AcceptLoop()
}

func (sv *server) AcceptLoop() {
	for {
		conn, err := sv.ln.Accept()
		if err != nil {
			break
		}
		fmt.Printf("client connected at: %v\n", conn.RemoteAddr().String())
		buffer := fmt.Appendf(nil, "connected to %v\n", conn.LocalAddr().String())
		conn.Write(buffer)
		// var wg sync.WaitGroup
		// start := time.Now()
		sv.wg.Add(1) // maybe use one global waitgroup
		go sv.HandleConn(conn)
		// fmt.Printf("[%v] waiting for handleconn to finish\n", time.Now().Format(time.TimeOnly))
		// wg.Wait()
		// fmt.Printf("handle connection finished in %v\n", time.Since(start))
		// fmt.Printf("[%v] waiting for handle conn finished\n", time.Now().Format(time.TimeOnly))
	}
}

func (sv *server) HandleConn(conn net.Conn) {
	defer conn.Close()
	defer sv.wg.Done()
	ctx, done := context.WithCancel(context.Background())
	defer done()
	go func() { // this wont close after conn.close
		select {

		case <-sv.ctx.Done():
			fmt.Println("closing connection")
			conn.Close()
		case <-ctx.Done():
		}
	}()
	sv.ReadLoop(conn)
}

func (sv *server) ReadLoop(conn net.Conn) {
	reader := bufio.NewReader(conn)
	workers := 10
	jobs := make(chan string, workers)
	// var fanoutWg sync.WaitGroup
	// fanoutWg.Add(1)
	// wg.Add(1)
	// var wg sync.WaitGroup
	var sendJobWg sync.WaitGroup
	sendJobWg.Add(1)
	// sv.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		// defer sv.wg.Done()
		defer wg.Done()
		defer close(jobs)
		for {
			line, _, err := reader.ReadLine() // back pressure, i think
			if err != nil {
				log.Printf("client %v disconnected %v\n", conn.RemoteAddr().String(), err) // so program exit after the client leaves, to track time
				return
			}
			// wg.Go(func() {
			jobs <- string(line)
			// })
			// wg.Wait()
		}
	}(&sendJobWg)
	// fmt.Printf("[%v] waiting for fanout to finish...\n", time.Now().Format(time.TimeOnly))
	sv.fanOut(jobs, workers)
	// fmt.Printf("[%v] waiting for fanout finished.\n", time.Now().Format(time.TimeOnly))

	fmt.Printf("[%v] waiting for sending jobs...\n", time.Now().Format(time.TimeOnly))
	sendJobWg.Wait()
	fmt.Printf("[%v] waiting for sending jobs finished.\n", time.Now().Format(time.TimeOnly))

}

func (sv *server) process(job string) ([]byte, error) {
	level := strings.Split(job, " ")[3]
	log := pkg.LogType{
		Level: level,
	}
	val, err := json.Marshal(log)
	time.Sleep(3 * time.Second) // simulate processing
	return val, err
}

func (sv *server) fanOut(jobs <-chan string, workers int) {
	var workerWg sync.WaitGroup
	for worker := range workers {
		workerWg.Add(1)
		sv.wg.Add(1)
		go func(worker int, workerWg *sync.WaitGroup) {
			defer workerWg.Done()
			defer sv.wg.Done()
			for job := range jobs {
				// start := time.Now()
				fmt.Printf("[%v] job \"%v\" being done by the worker %v...\n", time.Now().Format(time.TimeOnly), job, worker)
				log, err := sv.process(job)
				if err != nil {
					fmt.Printf("failed to process log %v", err)
				}
				sv.rdb.LPush(context.Background(), "job", string(log))
				fmt.Printf("[%v] worker %v finished the job %v.\n", time.Now().Format(time.TimeOnly), worker, job)
			}
		}(worker, &workerWg)
	}
	fmt.Printf("[%v] waiting for all workers to finish...\n", time.Now().Format(time.TimeOnly))
	workerWg.Wait() // waiting for fanout to finish and fanout() is syncronous so this is blocking, job is never being sent deadlock ?
	fmt.Printf("[%v] waiting for all workers finished\n", time.Now().Format(time.TimeOnly))
}
