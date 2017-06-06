package schema

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func drainAndCheckEqual(t *testing.T, req *http.Request, body []byte) {
	drained, err := drainHTTPRequestBody(req)
	if err != nil {
		t.Errorf("error during drain: %v", err)
		return
	}

	if bytes.Compare(drained, body) != 0 {
		t.Errorf("drained unexpected body: got %v want %v", drained, body)
		return
	}
}

func TestDrainHTTPRequestBody(t *testing.T) {
	var someTestBody []byte = []byte(`SOME TEST BODY`)
	req, err := http.NewRequest("POST", "/", bytes.NewReader(someTestBody))
	if err != nil {
		t.Errorf("cannot create new request: %v", err)
		return
	}

	// drain shouldn't affect further reads, do it multiple times and check result
	drainAndCheckEqual(t, req, someTestBody)
	drainAndCheckEqual(t, req, someTestBody)
	drainAndCheckEqual(t, req, someTestBody)
	drainAndCheckEqual(t, req, someTestBody)
}

func TestDrainHTTPRequestBodyNil(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Errorf("cannot create new request: %v", err)
		return
	}

	body, err := drainHTTPRequestBody(req)
	if err == nil || err != io.EOF {
		t.Errorf("Error should be %v, received %v", io.EOF, err)
	}
	if body != nil {
		t.Errorf("Drained body should be nil, received %v", body)
	}
}
