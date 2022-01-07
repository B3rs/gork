package workers

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
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
		for i := 0; i < info.instances; i++ {
			wg.Add(1)
			go w.workerRoutine(ctx, name, info.worker, errChan, wg)
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

func (w *WorkerPool) workerRoutine(ctx context.Context, queueName string, worker Worker, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	newWorker(w.db, queueName, worker, w.pollSleepInterval).startWorkLoop(ctx, errChan)
}
