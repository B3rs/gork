package workers

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

func NewWorkerPool(db *sql.DB, pollSleepInterval time.Duration) *WorkerPool {
	return &WorkerPool{
		db:                db,
		register:          newRegister(),
		pollSleepInterval: pollSleepInterval,
	}
}

type WorkerPool struct {
	register
	db *sql.DB

	pollSleepInterval time.Duration

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

	for name, info := range w.register.getWorkers() {

		// need to spawn a goroutine for polling the db
		poller := newPoller(w.db, name, w.pollSleepInterval)
		wg.Add(1)
		go pollerRoutine(ctx, poller, errChan, wg)

		for i := 0; i < info.instances; i++ {
			wg.Add(1)
			go workerRoutine(ctx, info.worker, poller.jobs(), errChan, wg)
		}
	}

	// wait for workers and pollers to stop
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

func pollerRoutine(ctx context.Context, poller poller, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	defer poller.stop()
	defer fmt.Println("Shutting down poller...")

	publishErrorAndRestart(errChan, func() error {
		return poller.start(ctx)
	})
}
func workerRoutine(ctx context.Context, worker Worker, c <-chan jobs.Job, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	execute(ctx, worker, c, errChan)
}

func publishErrorAndRestart(errChan chan<- error, f func() error) {
	for {
		err := f()
		if err == nil {
			return
		}
		errChan <- err
	}
}
