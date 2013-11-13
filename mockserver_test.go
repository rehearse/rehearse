package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type processedResponse struct {
	statusCode int
	header     http.Header
	body       []byte
}

func mustEncodeStub(t *testing.T, stub Stub) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(stub); err != nil {
		t.Fatalf("Unable to JSON encode: %#v", stub)
	}
	return
}

func mustPost(t *testing.T, url string, body io.Reader) (pr *processedResponse) {
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		t.Fatalf("Unable to POST: %v", err)
	}

	pr = &processedResponse{statusCode: resp.StatusCode, header: resp.Header}
	pr.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()

	return
}

func mustGet(t *testing.T, url string) (pr *processedResponse) {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	pr = &processedResponse{statusCode: resp.StatusCode, header: resp.Header}
	pr.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()

	return
}

func TestStubEndpoint(t *testing.T) {
	handler := NewStubHandler()
	ts := httptest.NewServer(handler)

	stub := Stub{Method: "GET", Path: "/foo", Body: `{"foo":"bar"}`}
	postBody := mustEncodeStub(t, stub)
	mustPost(t, ts.URL+"/stubs", postBody)

	resp := mustGet(t, ts.URL+"/foo")
	if resp.statusCode != 200 {
		t.Errorf("Status was incorrect, expected 200, got: %v", resp.statusCode)
	}

	if string(resp.body) != `{"foo":"bar"}` {
		t.Errorf("Unexpected response body: %s", string(resp.body))
	}
}

func TestStubEndpointWithMultipleMethods(t *testing.T) {
	handler := NewStubHandler()
	ts := httptest.NewServer(handler)

	stub := Stub{Method: "GET", Path: "/foo", Body: `get response`}
	postBody := mustEncodeStub(t, stub)
	mustPost(t, ts.URL+"/stubs", postBody)

	stub = Stub{Method: "POST", Path: "/foo", Body: `post response`}
	postBody = mustEncodeStub(t, stub)
	mustPost(t, ts.URL+"/stubs", postBody)

	resp := mustGet(t, ts.URL+"/foo")
	if resp.statusCode != 200 {
		t.Errorf("Status was incorrect, expected 200, got: %v", resp.statusCode)
	}

	if string(resp.body) != `get response` {
		t.Errorf("Unexpected response body: %s", string(resp.body))
	}

	resp = mustPost(t, ts.URL+"/foo", &bytes.Buffer{})
	if resp.statusCode != 200 {
		t.Errorf("Status was incorrect, expected 200, got: %v", resp.statusCode)
	}

	if string(resp.body) != `post response` {
		t.Errorf("Unexpected response body: %s", string(resp.body))
	}
}

func TestGetStubs(t *testing.T) {
	var stubs []Stub
	handler := NewStubHandler()
	ts := httptest.NewServer(handler)

	stub := Stub{Method: "GET", Path: "/foo", Body: `{"foo":"bar"}`}
	postBody := mustEncodeStub(t, stub)
	mustPost(t, ts.URL+"/stubs", postBody)
	resp := mustGet(t, ts.URL+"/stubs")
	decoder := json.NewDecoder(bytes.NewBuffer(resp.body))
	decoder.Decode(&stubs)
	if len(stubs) != 1 {
		t.Errorf("Expected one stub, got %d.", len(stubs))
	}
}

func TestClear(t *testing.T) {
	var stubs []Stub
	handler := NewStubHandler()
	ts := httptest.NewServer(handler)

	stub := Stub{Method: "GET", Path: "/foo", Body: `{"foo":"bar"}`}
	postBody := mustEncodeStub(t, stub)
	mustPost(t, ts.URL+"/stubs", postBody)
	resp := mustGet(t, ts.URL+"/stubs")
	decoder := json.NewDecoder(bytes.NewBuffer(resp.body))
	decoder.Decode(&stubs)
	if len(stubs) != 1 {
		t.Errorf("Expected one stub, got %d.", len(stubs))
	}

	client := http.Client{}
	req, err := http.NewRequest("DELETE", ts.URL+"/stubs", bytes.NewBuffer([]byte("")))
	client.Do(req)
	if err != nil {
		t.Error(err)
	}
	resp = mustGet(t, ts.URL+"/stubs")
	decoder = json.NewDecoder(bytes.NewBuffer(resp.body))
	decoder.Decode(&stubs)
	if len(stubs) != 0 {
		t.Errorf("Expected zero stubs, got %d.", len(stubs))
	}
}

func TestStubMalformedRequest(t *testing.T) {
	handler := NewStubHandler()
	ts := httptest.NewServer(handler)

	postBody := bytes.NewBuffer([]byte("not JSON"))
	resp, err := http.Post(ts.URL+"/stubs", "application/json", postBody)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Status was incorrect, expected 400, got: %v", resp.StatusCode)
	}
}

func TestUnhandledStub(t *testing.T) {
	handler := NewStubHandler()
	ts := httptest.NewServer(handler)

	resp, err := http.Get(ts.URL + "/bar")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("Status was incorrect, expected 404, got: %v", resp.StatusCode)
	}
}
