// Package data is responsible for encapsulating database specific code.
package data

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ErrNoRowsAfected = errors.New("no rows was affected")

type Config struct {
	Address  string `config:"DB_URL"`
	DataBase string `config:"DB_NAME"`
	User     string `config:"DB_USER"`
	Password string `config:"DB_PASSWORD"`
}

// Idiomatic way to pass default values in golang is in empty fields
// This function is because i don't want have default value ifs
// all over code
func populateConfig(c Config) Config {

	// database is mandatory
	if c.DataBase == "" {
		c.DataBase = "wp"
	}

	// if password specified but not user - asume root
	// mysql driver expects only user or both
	if c.Password != "" && c.User == "" {
		c.User = "root"
	}

	return c
}

// Data model with relevant for us geters and setters.
type Model interface {

	// Puts task into database, returns it ID
	PutTask(t Task) (ID uint, err error)

	// Lists tasks
	ListTasks() ([]Task, error)

	// Update task
	UpdateTask(t Task) error

	// Delete task
	DeleteTask(id uint) error

	// Put single entry
	PutEntry(e Entry) (uint, error)

	// List entries of single task
	ListEntries(id uint) ([]Entry, error)

	// Closes connection to database.
	Close() error

	// Verifies connection to database, reestablish it if needed.
	Ping() error
}

type data struct {
	config Config
	db     *gorm.DB
}

// Initializes datastore model.
func InitModel(c Config) (m Model, err error) {
	d := &data{
		config: populateConfig(c),
	}

	dsn := "/" + d.config.DataBase

	if d.config.Address != "" {
		dsn = "tcp(" + d.config.Address + ")" + dsn
	}

	if d.config.User != "" {
		u := d.config.User
		if d.config.Password != "" {
			u += ":" + d.config.Password
		}
		dsn = u + "@" + dsn
	}

	d.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	err = d.db.AutoMigrate(&Task{})
	if err != nil {
		return
	}

	err = d.db.AutoMigrate(&Entry{})
	if err != nil {
		return
	}

	return d, nil
}

func (d *data) Close() error {
	dbConn, err := d.db.DB()
	if err != nil {
		return err
	}
	return dbConn.Close()
}

func (d *data) Ping() error {
	dbConn, err := d.db.DB()
	if err != nil {
		return err
	}
	return dbConn.Ping()
}
