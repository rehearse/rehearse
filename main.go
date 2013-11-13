package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type StubHandler struct {
	stubs           map[string]Stub
	fallbackHandler http.Handler
	stubMutex       sync.Mutex
}

type Stub struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
}

func (h *StubHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Always lock the stubs map for the duration of the request.
	// Note that this effectively serializes all calls through StubHandler.
	// As this is only used for testing this is deemed acceptable.
	h.stubMutex.Lock()
	defer h.stubMutex.Unlock()

	if req.URL.Path == "/stubs" {
		switch req.Method {
		case "POST":
			h.createStubHandler(w, req)
		case "GET":
			h.listStubsHandler(w, req)
		case "DELETE":
			h.deleteStubsHandler(w, req)
		}
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
		_, err = fmt.Fprintf(w, "Invalid JSON: %v", err)
		if err != nil {
			log.Printf("Could not send response to client due to: %v", err)
		}
	}

	h.stubs[stub.Method+stub.Path] = stub
}

func (h *StubHandler) returnStubHandler(w http.ResponseWriter, req *http.Request) {
	if stub, ok := h.stubs[req.Method+req.URL.Path]; ok {
		_, err := io.WriteString(w, stub.Body)
		if err != nil {
			log.Printf("Could not send response to client due to: %v", err)
		}
	} else {
		if h.fallbackHandler == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			h.fallbackHandler.ServeHTTP(w, req)
		}
	}
}

func (h *StubHandler) listStubsHandler(w http.ResponseWriter, req *http.Request) {
	var stubs []Stub
	encoder := json.NewEncoder(w)
	for _, v := range h.stubs {
		stubs = append(stubs, v)
	}

	encoder.Encode(stubs)
}

func (h *StubHandler) deleteStubsHandler(w http.ResponseWriter, req *http.Request) {
	h.stubs = make(map[string]Stub)
}

func NewStubHandler() *StubHandler {
	var s StubHandler
	s.stubs = make(map[string]Stub)
	return &s
}

func main() {
}
