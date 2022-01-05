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

	ctx    context.Context
	cancel context.CancelFunc
}

func execute(ctx context.Context, worker Worker, jobsChan <-chan jobs.Job, errChan chan<- error) {
	for {
		select {
		case job, ok := <-jobsChan:
			if !ok {
				return
			}

			// are we passing the correct context? do we need to create a new context with some info or timeout?
			err := work(ctx, worker, job)
			if err != nil {
				errChan <- err
				continue
			}
		case <-ctx.Done():
			fmt.Printf("exiting worker. Error detail: %v\n", ctx.Err())
			return
		}
	}
}

func (w WorkerPool) spawnWorkers(ctx context.Context) {
	wg := &sync.WaitGroup{}

	errChan := make(chan error)
	for name, info := range w.register.GetWorkers() {

		fmt.Println("Spawning workers:", name, info)

		// need to spawn a goroutine for polling the db
		poller := newPoller(w.db, name, w.pollSleepInterval)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer poller.stop()

			publishErrorAndRestart(errChan, func() error {
				return poller.start(ctx)
			})
		}()

		for i := 0; i < info.instances; i++ {
			wg.Add(1)
			go func(worker Worker, c <-chan jobs.Job) {
				defer wg.Done()
				execute(ctx, worker, c, errChan)
			}(info.worker, poller.jobs())
		}
	}

	errwg := &sync.WaitGroup{}
	errwg.Add(1)
	go func() {
		defer errwg.Done()
		wg.Wait()
		defer close(errChan)
	}()

	errwg.Add(1)
	go func() {
		defer errwg.Done()
		for err := range errChan {
			fmt.Printf("Error: %v\n", err)
		}
	}()
	errwg.Wait()
}

// Start the WorkerPool
func (w *WorkerPool) Start() {

	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.spawnWorkers(w.ctx)

}

// Stop the WorkerPool
func (w WorkerPool) Stop() {
	w.cancel()
}

func publishErrorAndRestart(errChan chan<- error, f func() error) {
	for {
		err := f()
		if err == context.Canceled {
			return
		}
		errChan <- err
	}
}
