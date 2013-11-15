package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	VERSION = "0.2.0"
)

var options struct {
	path         string
	address      string
	port         int
	stubsPath    string
	printVersion bool
}

func main() {
	var err error

	flag.StringVar(&options.path, "path", "", "Optional path to serve static files from.")
	flag.StringVar(&options.address, "address", "127.0.0.1", "Address to listen on.")
	flag.IntVar(&options.port, "port", 3333, "Port to listen on.")
	flag.StringVar(&options.stubsPath, "stubs", "", "Optional path to JSON file to preload stubs.")
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

	if options.stubsPath != "" {
		func() {
			file, err := os.Open(options.stubsPath)
			if err != nil {
				log.Fatalf("Unable to open %s: %v", options.stubsPath, err)
			}
			defer file.Close()

			if err := stubHandler.load(file); err != nil {
				log.Fatalf("Unable to load %s: %v", options.stubsPath, err)
			}
		}()
	}

	http.Handle("/", stubHandler)

	log.Printf("Starting rehearse on %s:%d", options.address, options.port)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", options.address, options.port), stubHandler)
	if err != nil {
		log.Fatalf("Unable to start server: %v\n", err)
	}
}
