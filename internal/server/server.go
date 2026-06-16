package server

import (
	"fmt"
	"log"
	"net"
)

type server struct {
	ln net.Listener
}

func NewServer() server {
	return server{}
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

func (*server) ReadLoop(conn net.Conn, buffer []byte) {
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("client %v disconnected %v\n", conn.RemoteAddr().String(), err)
			return
		}
		fmt.Println(string(buffer[:n]))
	}
}
