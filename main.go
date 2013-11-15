package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	VERSION = "0.1.0"
)

var options struct {
	path         string
	address      string
	port         int
	printVersion bool
}

func main() {
	flag.StringVar(&options.path, "path", "", "Optional path to serve static files from.")
	flag.StringVar(&options.address, "address", "127.0.0.1", "Address to listen on.")
	flag.IntVar(&options.port, "port", 3333, "Port to listen on.")
	flag.BoolVar(&options.printVersion, "version", false, "Print version and exit.")
	flag.Parse()

	if options.printVersion {
		fmt.Printf("rehearse v%s\n", VERSION)
		os.Exit(0)
	}

	stubHandler := NewStubHandler()
	if options.path != "" {
		stubHandler.fallbackHandler = http.FileServer(http.Dir(options.path))
	}
	http.Handle("/", stubHandler)

	log.Printf("Starting rehearse on %s:%d", options.address, options.port)

	var err error
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", options.address, options.port), stubHandler)
	if err != nil {
		log.Fatalf("Unable to start server: %v\n", err)
	}
}
