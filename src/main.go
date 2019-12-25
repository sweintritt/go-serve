package main

import (
	"flag"
)

func main() {
	server := NewServer()

	port := ""
	flag.StringVar(&port, "p", server.Port, "Service port")
	flag.StringVar(&port, "port", server.Port, "Service port")
	flag.Parse()

	if port != "" {
		server.Port = port
	}

	server.Start()
}
