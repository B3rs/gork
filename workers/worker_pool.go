package workers

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

type poolConfig struct {
	schedulerSleepInterval time.Duration
	reaperInterval         time.Duration
}

func NewWorkerPool(db *sql.DB, opts ...PoolOptionFunc) *WorkerPool {
	po := poolConfig{}

	options := append(defaultPoolOptions, opts...)
	for _, opt := range options {
		po = opt(po)
	}

	return &WorkerPool{
		db:       db,
		register: newRegister(),
		options:  po,
	}
}

type WorkerPool struct {
	register
	db      *sql.DB
	options poolConfig

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
	go errorRoutine(errChan, errwg)

	wg := &sync.WaitGroup{}

	for name, config := range w.register.getWorkers() {
		q := jobs.NewQueue(w.db, name)

		// worker routines
		for i := 0; i < config.instances; i++ {
			wg.Add(1)
			go workerRoutine(ctx, q, config.worker, w.options.schedulerSleepInterval, errChan, wg)
		}

		// reaper routine
		wg.Add(1)
		go reaperRoutine(ctx, q, w.options.reaperInterval, config.timeout, errChan, wg)
	}

	// wait for workers and reapers to stop
	wg.Wait()

	// wait for error routine to stop
	close(errChan)
	errwg.Wait()
}

func errorRoutine(errChan <-chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range errChan {
		fmt.Printf("Error: %v\n", err)
	}
}

func workerRoutine(ctx context.Context, queue queue, worker Worker, sleepInterval time.Duration, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	s := newScheduler(queue, worker, sleepInterval)
	s.Run(ctx, errChan)
}

func reaperRoutine(ctx context.Context, queue requeuer, reaperInterval time.Duration, jobTimeout time.Duration, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	r := newReaper(queue, reaperInterval, jobTimeout)
	r.Run(ctx, errChan)
}
