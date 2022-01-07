package workers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

// WorkerFunc is a function that can be used as a worker.
type WorkerFunc func(ctx context.Context, job jobs.Job) (interface{}, error)

// Worker is the interface that must be implemented by workers.
type Worker interface {
	Execute(ctx context.Context, job jobs.Job) (interface{}, error)
}

type store interface {
	Begin() (*jobs.Tx, error)
}

func newWorker(db *sql.DB, queueName string, w Worker, sleepInterval time.Duration) *worker {
	return &worker{
		store:         jobs.NewStore(db),
		queueName:     queueName,
		worker:        w,
		sleepInterval: sleepInterval,
	}
}

type worker struct {
	queueName     string
	worker        Worker
	sleepInterval time.Duration

	store store
}

// Work runs the worker function for the given job.
func (w *worker) work(ctx context.Context, tx *jobs.Tx, job *jobs.Job) error {

	// explicity copy the job to avoid user provided function to modify the job
	if res, err := w.worker.Execute(ctx, *job); err != nil {
		job.SetStatus(jobs.StatusFailed)
		job.SetLastError(err)
	} else {
		job.SetStatus(jobs.StatusCompleted)
		if err := job.SetResult(res); err != nil {
			job.SetLastError(err)
		}
	}
	return tx.Update(ctx, job)
}

func (w *worker) startWorkLoop(ctx context.Context, errChan chan<- error) {
	defer fmt.Println("Shutting down worker...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			tx, err := w.store.Begin()
			if err != nil {
				errChan <- err
				continue
			}

			job, err := tx.GetAndLock(ctx, w.queueName)
			if err == jobs.ErrJobNotFound {
				if err := tx.Commit(); err != nil {
					errChan <- err
				}
				time.Sleep(w.sleepInterval)
				continue
			}
			if err != nil {
				errChan <- err
				continue
			}

			if err := w.work(context.Background(), tx, job); err != nil {
				errChan <- err
				continue
			}

			if err := tx.Commit(); err != nil {
				errChan <- err
			}
		}
	}
}
