// Http handlers reside here. Logicaly very simple package.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SQLek/wp-interview/model"
	"github.com/SQLek/wp-interview/worker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	Endpoint        string
	Timeout         time.Duration
	MaxHeader       int
	MaxBody         int
	WorkerMode      bool
	DefaultInterval uint
}

// Idiomatic way to pass default values in golang is in empty fields
// This function is because i don't want have default value ifs
// all over code
func populateConfig(c Config) Config {
	if c.MaxBody <= 0 {
		log.Println("MAX_BODY unspecified. Defaulting to 5kb")
		c.MaxBody = 5 * 1024
	}
	if c.Endpoint == "" {
		log.Println("ENDPOINT unspecified. Defaulting to ':8080'")
		c.Endpoint = ":8080"
	}
	if c.DefaultInterval == 0 {
		log.Println("DEFAULT_INTERVAL unspecified. Defaulting to '60s'")
		c.DefaultInterval = 60
	}
	return c
}

// contekstual arguments for handlers
type server struct {
	config Config
	model  model.Model
	sche   worker.Scheduler
}

// Makes http.Server with chi router and handlers
func MakeServer(c Config, m model.Model, w worker.Scheduler) *http.Server {

	s := &server{
		config: populateConfig(c),
		model:  m,
		sche:   w,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/api/urls", s.listTasks)
	r.Post("/api/urls", s.createTask)
	r.Patch("/api/urls/{taskID}", s.updateTask)
	r.Delete("/api/urls/{taskID}", s.deleteTask)
	r.Get("/api/urls/{taskID}/history", s.inspectTask)

	return &http.Server{
		Addr:           s.config.Endpoint,
		Handler:        r,
		ReadTimeout:    s.config.Timeout,
		WriteTimeout:   s.config.Timeout,
		MaxHeaderBytes: s.config.MaxHeader,
	}
}

// PUT /api/urls handler
func (s *server) createTask(w http.ResponseWriter, r *http.Request) {
	reader := http.MaxBytesReader(w, r.Body, int64(s.config.MaxBody))
	decoder := json.NewDecoder(reader)

	task := model.Task{}
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	if task.URL == "" {
		http.Error(w, "Bad Request", 400)
		return
	}
	if task.Interval == 0 {
		task.Interval = s.config.DefaultInterval
	}

	i, err := s.sche.SpawnWorker(r.Context(), task)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, `{"id": %d}`, i)
	if err != nil {
		// only posibility this would fail is break in connection
		log.Println(err.Error())
	}

}

// PATCH /api/urls/{id} handler
func (s *server) updateTask(w http.ResponseWriter, r *http.Request) {
	reader := http.MaxBytesReader(w, r.Body, int64(s.config.MaxBody))
	decoder := json.NewDecoder(reader)

	task := model.Task{}
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	if task.URL == "" {
		http.Error(w, "Bad Request", 400)
		return
	}
	if task.Interval == 0 {
		task.Interval = s.config.DefaultInterval
	}

	tid, err := strconv.ParseUint(chi.URLParam(r, "taskID"), 10, 32)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	task.ID = uint(tid)

	err = s.sche.UpdateWorker(r.Context(), task)
	if err != nil {
		if errors.Is(err, model.ErrNoRowsAfected) {
			http.Error(w, "Not Found", 404)
		} else {
			http.Error(w, "Internal server error", 500)
		}
		return
	}

}

// DELETE /api/urls/{id} handler
func (s *server) deleteTask(w http.ResponseWriter, r *http.Request) {

	tid, err := strconv.ParseUint(chi.URLParam(r, "taskID"), 10, 32)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	err = s.sche.KillWorker(r.Context(), uint(tid))
	if err != nil {
		if errors.Is(err, model.ErrNoRowsAfected) {
			http.Error(w, "Not Found", 404)
		} else {
			http.Error(w, "Internal server error", 500)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/urls/{id}/history
func (s *server) inspectTask(w http.ResponseWriter, r *http.Request) {

	tid, err := strconv.ParseUint(chi.URLParam(r, "taskID"), 10, 32)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	entries, err := s.model.ListEntries(uint(tid))
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&entries)
	if err != nil {
		http.Error(w, "Internal server error", 500)
	}

}

// GET /api/urls handler
func (s *server) listTasks(w http.ResponseWriter, r *http.Request) {

	tasks, err := s.model.ListTasks()
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&tasks)
	if err != nil {
		http.Error(w, "Internal server error", 500)
	}
}
