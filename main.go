package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	VERSION = "0.1"
)

var config struct {
	path         string
	address      string
	port         int
	printVersion bool
}

func main() {
	flag.StringVar(&config.path, "path", "", "Optional path to serve static files from.")
	flag.StringVar(&config.address, "address", "127.0.0.1", "Address to listen on.")
	flag.IntVar(&config.port, "port", 3333, "Port to listen on.")
	flag.BoolVar(&config.printVersion, "version", false, "Print version and exit.")
	flag.Parse()

	if config.printVersion {
		fmt.Printf("rehearse v%s\n", VERSION)
		os.Exit(0)
	}

	stubHandler := NewStubHandler()
	if config.path != "" {
		stubHandler.fallbackHandler = http.FileServer(http.Dir(config.path))
	}
	http.Handle("/", stubHandler)

	var err error
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.address, config.port), stubHandler)
	if err != nil {
		log.Fatalf("Unable to start server: %v\n", err)
	}
}
