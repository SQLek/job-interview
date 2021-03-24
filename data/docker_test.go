package data

import (
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
)

var testModel Model = nil

func TestMain(m *testing.M) {
	// Why mock database driver or orm when You can just run real thing
	// Awesome library! https://github.com/ory/dockertest
	// modified version of example from His github

	config := Config{
		DataBase: "test_db",
		Password: "test_password",
		// No user, code should asume root user
		Address: "localhost",
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8", []string{
		"MYSQL_ROOT_PASSWORD=" + config.Password,
		"MYSQL_DATABASE=" + config.DataBase})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// outside port of database can be different
	config.Address += ":" + resource.GetPort("3306/tcp")

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() (err error) {

		testModel, err = InitModel(config)
		if err != nil {
			return err
		}
		return testModel.Ping()
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
