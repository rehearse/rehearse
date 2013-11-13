package main

import (
	"encoding/json"
	"io"
	"log"
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
		h.createStubHandler(w, req)
	} else {
		h.returnStubHandler(w, req)
	}
}

func (h *StubHandler) createStubHandler(w http.ResponseWriter, req *http.Request) {
	var stub Stub
	jsonDecoder := json.NewDecoder(req.Body)
	err := jsonDecoder.Decode(&stub)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	h.stubs[stub.Path] = stub
}

func (h *StubHandler) returnStubHandler(w http.ResponseWriter, req *http.Request) {
	if stub, ok := h.stubs[req.URL.Path]; ok {
		if stub.Method == req.Method {
			_, err := io.WriteString(w, stub.Body)
			if err != nil {
				log.Printf("Could not send response to client due to: %v", err)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func NewStubHandler() *StubHandler {
	var s StubHandler
	s.stubs = make(map[string]Stub)
	return &s
}

func main() {
}
