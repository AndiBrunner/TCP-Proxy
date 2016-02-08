package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/nseveryns/proxy/proxy"
)

var (
	bind    = flag.String("b", ":9999", "Address to bind on")
	forward = flag.String("f", "localhost:1000", "IP to forward connections to")
)

func main() {
	flag.Parse()
	p := proxy.New(*bind, *forward)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		p.Close()
		os.Exit(1)
	}()

	p.Start()
}
