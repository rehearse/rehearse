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

func mustEncodeStub(t *testing.T, stub Stub) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(stub); err != nil {
		t.Fatalf("Unable to JSON encode: %#v", stub)
	}
	return
}

func mustPost(t *testing.T, url string, body io.Reader) {
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		t.Fatalf("Unable to POST: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Post failure - Status was incorrect, expected 200, got: %v", resp.Status)
	}

	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("Closing response body failed: %v", err)
	}
}

func mustGet(t *testing.T, url string) (resp *http.Response) {
	var err error
	resp, err = http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Status was incorrect, expected 200, got: %v", resp.Status)
	}

	return
}

func TestStubEndpoint(t *testing.T) {
	handler := NewStubHandler()
	ts := httptest.NewServer(handler)

	stub := Stub{Method: "GET", Path: "/foo", Body: `{"foo":"bar"}`}
	postBody := mustEncodeStub(t, stub)
	mustPost(t, ts.URL+"/stub", postBody)

	resp := mustGet(t, ts.URL+"/foo")

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if string(respBody) != `{"foo":"bar"}` {
		t.Errorf("Unexpected response body: %s", string(respBody))
	}
}
