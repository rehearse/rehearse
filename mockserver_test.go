package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStubEndpoint(t *testing.T) {
	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	handler := NewStubHandler()
	ts := httptest.NewServer(handler)
	stub := Stub{Method: "GET", Path: "/foo", Body: `{"foo":"bar"}`}
	encoder.Encode(stub)
	resp, err := http.Post(ts.URL+"/stub", "text/json", body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Status was incorrect, expected 200, got: %v", resp.Status)
	}

	resp.Body.Close()

	body.Reset()
	resp, err = http.Get(ts.URL + "/foo")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Status was incorrect, expected 200, got: %v", resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if string(respBody) != `{"foo":"bar"}` {
		t.Errorf("Unexpected response body: %s", string(respBody))
	}
}
