package workers

import (
	"context"
	"fmt"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

type WorkerFunc func(ctx context.Context, job jobs.Job) (interface{}, error)

type Worker interface {
	Execute(ctx context.Context, job jobs.Job) (interface{}, error)
}

// Work runs the worker function for the given job.
func work(ctx context.Context, worker Worker, job jobs.Job) error {
	res, err := worker.Execute(ctx, job)
	if err != nil {
		if err := job.SetStatus(ctx, jobs.StatusFailed); err != nil {
			return err
		}

		if err := job.SetLastError(ctx, err); err != nil {
			return err
		}
		return job.Commit()
	}

	if err := job.SetStatus(ctx, jobs.StatusCompleted); err != nil {
		return err
	}

	if err := job.SetResult(ctx, res); err != nil {
		return err
	}

	return job.Commit()
}

func execute(ctx context.Context, worker Worker, jobsChan <-chan jobs.Job, errChan chan<- error) {
	defer fmt.Println("Shutting down worker...")
	for {
		select {
		case job, ok := <-jobsChan:
			if !ok {
				return
			}

			if err := work(context.Background(), worker, job); err != nil {
				errChan <- err // here there will be db errors and timeouts, not job errors that will be saved on db
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}
