package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const urlPrefix = "http://localhost:8080"

func TestMalformedPost(t *testing.T) {
	assert := assert.New(t)

	payload := []byte("qwertyuiop")
	_, code, err := requestID("POST", urlPrefix+"/api/urls", bytes.NewReader(payload))
	assert.Nil(err)
	assert.Equal(code, 400, "400 - Bad request expected")

}

func TestMalformedPatch(t *testing.T) {
	assert := assert.New(t)

	payload := []byte("qwertyuiop")
	_, code, err := requestID("PATCH", urlPrefix+"/api/urls/1", bytes.NewReader(payload))
	assert.Nil(err)
	assert.Equal(code, 400, "400 - Bad request expected")

}

func TestOversizedPost(t *testing.T) {
	assert := assert.New(t)

	payload, err := prepareBigPayload()
	assert.Nil(err)
	assert.Greater(len(payload), 5*1024, "Payload generation failed!")

	_, code, err := requestID("POST", urlPrefix+"/api/urls", bytes.NewReader(payload))
	assert.Nil(err)
	assert.Equal(code, 400, "400 - Bad request expected")

}

func TestOversizedPatch(t *testing.T) {
	assert := assert.New(t)

	payload, err := prepareBigPayload()
	assert.Nil(err)
	assert.Greater(len(payload), 5*1024, "Payload generation failed!")

	_, code, err := requestID("PATCH", urlPrefix+"/api/urls/1", bytes.NewReader(payload))
	assert.Nil(err)
	assert.Equal(code, 400, "400 - Bad request expected")

}

func TestBadIdDelete(t *testing.T) {
	assert := assert.New(t)

	// TODO add code for checking what ID is not in use
	_, code, err := requestID("DELETE", urlPrefix+"/api/urls/666", nil)
	assert.Nil(err)
	assert.Equal(404, code, "Not found expected")

}

func TestInsertAndDeleteShort(t *testing.T) {
	assert := assert.New(t)

	payload := []byte(`{"url":"http://httpbin/range/15","interval":60}`)
	data, code, err := requestID("POST", urlPrefix+"/api/urls", bytes.NewReader(payload))
	assert.Nil(err)
	assert.Equal(201, code, "'Created' expected")

	ido := IDObject{}
	err = json.Unmarshal(data, &ido)
	assert.Nil(err)

	url := fmt.Sprintf("%s/api/urls/%d", urlPrefix, ido.ID)
	_, code, err = requestID("DELETE", url, nil)
	assert.Nil(err)
	assert.Equal(204, code, "'No content' expected")

}
