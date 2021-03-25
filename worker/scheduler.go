// Worker package is responsible for fetching external websites,
// and scheduling goroutines
package worker

import (
	"context"
	"sync"
	"time"

	"github.com/SQLek/wp-interview/model"
)

type Config struct {
	MinInterval  uint          `config:"MIN_INTERVAL"`
	FetchTimeout time.Duration `config:"FETCH_TIMEOUT"`
}

func populateConfig(c Config) Config {
	if c.MinInterval == 0 {
		c.MinInterval = 5
	}
	if c.FetchTimeout == 0 {
		c.FetchTimeout = 5 * time.Second
	}
	return c
}

// could be useful with testing
type workerFunc func(context.Context, model.Model, model.Task, time.Duration)

// Scheduler is responsible for spawning and managing goroutines
// that fetch data from external websites
type Scheduler interface {

	// Spawns worker and adds task do model
	SpawnWorker(context.Context, model.Task) (uint, error)

	// Updates worker and data entry in model
	// Currently timings can be little off.
	// This method kils goroutine and start new.
	// Some form of time kalculation could be done to alievate this
	UpdateWorker(context.Context, model.Task) error

	// Kils worker and inform model to remove task and entries
	KillWorker(context.Context, uint) error

	// Kils all workers
	KillAll()
}

type scheduler struct {
	config Config
	model  model.Model

	// synchronized map with cancel funtions of every worker
	mutex     *sync.Mutex
	cancelers map[uint]context.CancelFunc

	// to have option with testing
	worker workerFunc
}

func MakeSheduler(c Config, m model.Model) Scheduler {
	return &scheduler{
		config:    populateConfig(c),
		model:     m,
		mutex:     &sync.Mutex{},
		worker:    WorkOnRequests,
		cancelers: make(map[uint]context.CancelFunc),
	}
}

func (sched *scheduler) SpawnWorker(ctx context.Context, task model.Task) (uint, error) {
	sched.mutex.Lock()
	defer sched.mutex.Unlock()

	if task.Interval < sched.config.MinInterval {
		task.Interval = sched.config.MinInterval
	}

	id, err := sched.model.PutTask(task)
	if err != nil {
		return 0, err
	}
	// task is delivered by value
	// id isn't propagated here
	task.ID = id

	// ctx is from http server and will timeout worker quickly
	workerCtx, cancel := context.WithCancel(context.Background())

	sched.cancelers[id] = cancel
	go sched.worker(workerCtx, sched.model, task, sched.config.FetchTimeout)
	return id, nil
}

func (sched *scheduler) KillWorker(ctx context.Context, id uint) error {
	sched.mutex.Lock()
	defer sched.mutex.Unlock()

	canceler, exists := sched.cancelers[id]
	if !exists {
		// TODO maybe refactor this error type here?
		return model.ErrNoRowsAfected
	}

	canceler()
	return sched.model.DeleteTask(id)
}

func (sched *scheduler) UpdateWorker(ctx context.Context, task model.Task) error {
	sched.mutex.Lock()
	defer sched.mutex.Unlock()

	canceler, exists := sched.cancelers[task.ID]
	if !exists {
		// TODO maybe refactor this error type here?
		return model.ErrNoRowsAfected
	}

	// TODO we could check last time fetch was executed
	err := sched.model.UpdateTask(task)
	if err != nil {
		return err
	}
	canceler()

	workerCtx, cancel := context.WithCancel(context.Background())

	sched.cancelers[task.ID] = cancel
	go sched.worker(workerCtx, sched.model, task, sched.config.FetchTimeout)
	return nil
}

func (sched *scheduler) KillAll() {
	sched.mutex.Lock()
	defer sched.mutex.Unlock()

	for _, canceler := range sched.cancelers {
		canceler()
	}

	sched.cancelers = make(map[uint]context.CancelFunc)
}
