package workers

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

func NewWorker(db *sql.DB, queueName string, workFunc WorkFunc) Worker {
	return Worker{
		queue:    NewQueue(queueName, db),
		workFunc: workFunc,
	}
}

type WorkFunc func(ctx context.Context, job jobs.Job) (interface{}, error)

type Worker struct {
	queue    *Queue
	workFunc WorkFunc
	sleep    time.Duration
}

// Work does the work with an at least once semantic
// it does not return error if the job errored, only if there are db problems
func Work(ctx context.Context, job jobs.Job, workFunc WorkFunc) error {

	res, err := workFunc(ctx, job)
	if err != nil {
		if err := job.SetStatus(ctx, jobs.StatusFailed); err != nil {
			return err
		}

		if err := job.SetLastError(ctx, err); err != nil {
			return err
		}
		return nil
	}

	if err := job.SetStatus(ctx, jobs.StatusCompleted); err != nil {
		return err
	}

	if err := job.SetResult(ctx, res); err != nil {
		return err
	}

	return nil
}

func WorkLoop(id int, jobsChan chan jobs.Job, f WorkFunc) error {
	for job := range jobsChan {

		if err := Work(context.Background(), job, f); err != nil {
			_ = job.Commit()
			return err
		}

		// release lock by committing
		if err := job.Commit(); err != nil {
			return err
		}

	}
	return nil
}

// Start the workers
func (w Worker) Start(workersCount int) error {
	wg := sync.WaitGroup{}

	jobsChan := make(chan jobs.Job)
	errChan := make(chan error)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func(id int, f WorkFunc) {
			// keep doing work until the channel is closed
			for {
				if err := WorkLoop(id, jobsChan, f); err != nil {
					errChan <- err
					continue
				}
				break
			}

			wg.Done()
		}(i, w.workFunc)
	}

	wg.Add(1)
	go func() {
		for err := range errChan {
			fmt.Println("worker pool error", err)
		}
		wg.Done()
	}()

Loop:
	for {
		// Acquire job and lock it
		job, err := w.queue.AcquireJob()
		switch {
		case err == ErrJobNotFound:
			time.Sleep(w.sleep)
			continue Loop
		case err == ErrQueueIsClosed:
			close(jobsChan)
			errChan <- err
			close(errChan)
			break Loop
		case err != nil:
			errChan <- err
			continue Loop
		}

		jobsChan <- job
	}

	wg.Wait()

	return nil
}

// Start the workers
func (w Worker) Stop() {
	w.queue.Close()
}
