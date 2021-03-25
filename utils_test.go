package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// Helper methods and structs for integration testing

// struct that contains only uint ID
// it will be returned from:
// * POST /api/urls
// * PATCH /api/urls/{ID}
type IDObject struct {
	ID uint `json:"id"`
}

// Used in integration tests
//
// * POST /api/urls
// * PATCH /api/urls/{ID}
// * DELETE /api/urls/{ID}
func requestID(method, url string, body io.Reader) (result []byte, code int, err error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	rsp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	code = rsp.StatusCode
	result, err = ioutil.ReadAll(rsp.Body)

	return
}

// Prepares 5kb+ payload that is well formed json as well
func prepareBigPayload() ([]byte, error) {
	payload := make([]string, 0)
	for len(payload) < 128 {
		payload = append(payload, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	}
	return json.Marshal(payload)
}
