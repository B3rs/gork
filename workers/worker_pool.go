package workers

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/B3rs/gork/jobs"
)

type logger interface {
	Log(keyvals ...interface{}) error
}

func NewWorkerPool(db *sql.DB, opts ...PoolOptionFunc) *WorkerPool {
	w := &WorkerPool{
		db:       db,
		register: newRegister(),
	}

	options := append(defaultPoolOptions, opts...)
	for _, opt := range options {
		w = opt(w)
	}

	return w
}

type WorkerPool struct {
	register
	db     *sql.DB
	logger logger

	schedulerSleepInterval time.Duration
	reaperInterval         time.Duration

	shutdown func()
}

// Stop the WorkerPool
func (w WorkerPool) Stop() {
	w.shutdown()
}

// Start the WorkerPool
func (w *WorkerPool) Start() {

	ctx, cancel := context.WithCancel(context.Background())
	w.shutdown = cancel

	errChan := make(chan error)
	errwg := &sync.WaitGroup{}
	errwg.Add(1)
	go errorRoutine(errChan, w.logger, errwg)

	wg := &sync.WaitGroup{}

	for name, config := range w.register.getWorkers() {
		q := jobs.NewQueue(w.db, name)

		// worker routines
		for i := 0; i < config.instances; i++ {
			wg.Add(1)
			go workerRoutine(ctx, q, config.worker, w.schedulerSleepInterval, errChan, wg)
		}

		// reaper routine
		wg.Add(1)
		go reaperRoutine(ctx, q, w.reaperInterval, config.timeout, errChan, wg)
	}

	// wait for workers and reapers to stop
	wg.Wait()

	// wait for error routine to stop
	close(errChan)
	errwg.Wait()
}

func errorRoutine(errChan <-chan error, logger logger, wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range errChan {
		logger.Log("msg", fmt.Sprintf("Error: %v\n", err))
	}
}

func workerRoutine(ctx context.Context, queue Queue, worker Worker, sleepInterval time.Duration, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	runner := newRunner(worker, queue)
	s := newScheduler(queue, runner, sleepInterval)
	s.Run(ctx, errChan)
}

func reaperRoutine(ctx context.Context, queue Requeuer, reaperInterval time.Duration, jobTimeout time.Duration, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	r := newReaper(queue, reaperInterval, jobTimeout)
	r.Run(ctx, errChan)
}
