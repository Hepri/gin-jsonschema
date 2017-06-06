package schema

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// read http request body, and replace it with copy in order to allow further
func drainHTTPRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, io.EOF
	}

	// read body
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// replace request body with new one
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	return buf, nil
}
