package workers

import (
	"context"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

type queue interface {
	Dequeue(ctx context.Context) (*jobs.Job, error)
	Update(ctx context.Context, job *jobs.Job) error
}

func newScheduler(queue queue, w Worker, sleepInterval time.Duration) *scheduler {
	return &scheduler{
		queue:         queue,
		worker:        w,
		sleepInterval: sleepInterval,
	}
}

type scheduler struct {
	queue         queue
	worker        Worker
	sleepInterval time.Duration
}

func (s *scheduler) Run(ctx context.Context, errChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			job, err := s.queue.Dequeue(ctx)
			if err == jobs.ErrNoJobsAvailable {
				time.Sleep(s.sleepInterval)
				continue
			}
			if err != nil {
				errChan <- err
				continue
			}

			ex := newExecutor(s.worker, s.queue.Update)
			if err := ex.execute(context.Background(), job); err != nil {
				errChan <- err
				continue
			}
		}
	}
}
