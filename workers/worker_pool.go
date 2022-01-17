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
	db           *sql.DB
	errorHandler func(error)

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
	go errorRoutine(errChan, w.errorHandler, errwg)

	wg := &sync.WaitGroup{}

	for name, config := range w.register.getWorkers() {
		q := jobs.NewQueue(w.db, name)

		// worker routines
		for i := 0; i < config.instances; i++ {
			runner := newJobRunner(config.worker, q)
			s := newScheduler(q, runner, w.schedulerSleepInterval)

			wg.Add(1)
			go routine(ctx, s, errChan, wg)
		}

		// reaper routine
		r := newReaper(q, w.reaperInterval, config.timeout)

		wg.Add(1)
		go routine(ctx, r, errChan, wg)
	}

	// wait for workers and reapers to stop
	wg.Wait()

	// wait for error routine to stop
	close(errChan)
	errwg.Wait()
}

func errorRoutine(errChan <-chan error, errorHandler func(error), wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range errChan {
		errorHandler(err)
	}
}

type runner interface {
	Run(context.Context, chan<- error)
}

func routine(ctx context.Context, r runner, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	r.Run(ctx, errChan)
}

func defaultErrorHandler(err error) {
	log.Println("error in worker pool", "error", err)
}
