package workers

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/B3rs/gork/jobs"
)

func NewWorkerPool(db *sql.DB, opts ...PoolOptionFunc) *WorkerPool {
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	w := &WorkerPool{
		register: newRegister(),
		spawner:  newSpawner(ctx, errChan),
		errChan:  errChan,
		shutdown: cancel,
		queueFactory: func(name string) Queue {
			return jobs.NewQueue(db, name)
		},
	}

	options := append(defaultPoolOptions, opts...)
	for _, opt := range options {
		w = opt(w)
	}

	return w
}

type WorkerPool struct {
	register
	errorHandler func(error)
	errChan      chan error
	shutdown     func()
	spawner      Spawner

	queueFactory func(name string) Queue

	schedulerSleepInterval time.Duration
	reaperInterval         time.Duration
}

// Stop the WorkerPool
func (w WorkerPool) Stop() {
	w.spawner.Shutdown()
}

// Start the WorkerPool
func (w *WorkerPool) Start() {

	errwg := &sync.WaitGroup{}
	errwg.Add(1)
	go errorRoutine(w.errChan, w.errorHandler, errwg)

	for name, config := range w.register.getWorkers() {
		q := w.queueFactory(name)

		// worker routines
		for i := 0; i < config.instances; i++ {
			s := newDequeuer(q, config.worker, w.schedulerSleepInterval)
			w.spawner.Spawn(s)
		}

		// reaper routine
		r := newReaper(q, w.reaperInterval, config.timeout)
		w.spawner.Spawn(r)
	}

	// wait for workers and reapers to stop
	w.spawner.Wait()

	// wait for error routine to stop
	close(w.errChan)
	errwg.Wait()
}

func errorRoutine(errChan <-chan error, errorHandler func(error), wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range errChan {
		errorHandler(err)
	}
}

func defaultErrorHandler(err error) {
	log.Println("error in worker pool", "error", err)
}
