package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StubHandler struct {
	stubs map[string]Stub
}

type Stub struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
}

func (h *StubHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/stub" && req.Method == "POST" {
		var stub Stub
		jsonDecoder := json.NewDecoder(req.Body)
		// handle decode error
		err := jsonDecoder.Decode(&stub)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%#v\n", stub)
		h.stubs[stub.Path] = stub
	} else {
		fmt.Printf("%#v\n", h.stubs)
		if stub, ok := h.stubs[req.URL.Path]; ok {
			if stub.Method == req.Method {
				// TODO - handle errors
				_, _ = io.WriteString(w, stub.Body)
			}
		} else {
			fmt.Println("not found")
		}
	}
}

func NewStubHandler() *StubHandler {
	var s StubHandler
	s.stubs = make(map[string]Stub)
	return &s
}

func main() {
	fmt.Println("Hello, world!")
}
