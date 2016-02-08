package proxy

import (
	"io"
	"log"
	"net"
)

//Proxy instance containing the list of servers.
type Proxy struct {
	bind     string
	forward  string
	listener net.Listener
}

//New proxy insance
func New(bind, forward string) *Proxy {
	return &Proxy{bind: bind, forward: forward}
}

//Start the proxy server
func (p *Proxy) Start() {
	listener, err := net.Listen("tcp", p.bind)
	if err != nil {
		return
	}
	p.listener = listener
	log.Println("Started proxy bound on", p.bind)
	for {
		if conn, err := listener.Accept(); err == nil {
			p.handle(conn)
		} else {
			log.Fatal("Nothing to read in connection? ", err)
			p.Close()
		}
	}
}

//Close the proxy instance
func (p *Proxy) Close() {
	err := p.listener.Close()
	var resp string
	if err == nil {
		resp = "Success."
	} else {
		resp = "Failed."
	}
	log.Println("Closing proxy server:", resp)
}

func (p Proxy) handle(up net.Conn) {
	defer up.Close()
	down, err := net.Dial("tcp", p.forward)
	if err != nil {
		return
	}
	defer down.Close()
	pipe(up, down)
}

func pipe(a, b net.Conn) error {
	errors := make(chan error, 1)
	copy := func(write, read net.Conn) {
		_, err := io.Copy(write, read)
		errors <- err
	}
	go copy(a, b)
	go copy(b, a)
	return <-errors
}
