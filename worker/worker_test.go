package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	assert := assert.New(t)
	// GET request that should not timeout

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	r, err := request(ctx, testUrlPrefix+"/delay/4")

	assert.Nil(err)
	assert.Equal(200, r.Code, "We expect 200 status from 4s delay endpoint.")

}

func TestMeow(t *testing.T) {
	assert := assert.New(t)
	// Checking value that we guarante httpbin will return

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	r, err := request(ctx, testUrlPrefix+"/base64/TWVvdw==")

	assert.Nil(err)
	assert.Equal(200, r.Code, "We expect 200 status from 4s delay endpoint.")
	assert.Equal("Meow", r.Content, "We expect content to be 'Meow'")

}

func TestTimeout(t *testing.T) {
	assert := assert.New(t)
	// GET request that should timeout

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	_, err := request(ctx, testUrlPrefix+"/delay/6")

	if assert.Error(err, "We expect error") {
		assert.ErrorIs(err, context.DeadlineExceeded, "We expect deadline error")
	}

}
