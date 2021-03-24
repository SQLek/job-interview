package worker

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
)

var testUrlPrefix string

func TestMain(m *testing.M) {
	// Why mock database driver or orm when You can just run real thing
	// Awesome library! https://github.com/ory/dockertest
	// modified version of example from His github

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("kennethreitz/httpbin", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// outside port of httpbin can be different
	testUrlPrefix = "http://localhost:" + resource.GetPort("80/tcp")

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() (err error) {

		_, err = http.Get(testUrlPrefix + "/get")
		return err

	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
