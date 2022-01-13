package workers

import (
	"context"
	"time"

	"github.com/B3rs/gork/jobs"
)

func newScheduler(queue Queue, runner Runner, sleepInterval time.Duration) *scheduler {
	return &scheduler{
		queue:         queue,
		runner:        runner,
		sleepInterval: sleepInterval,
	}
}

type scheduler struct {
	queue         Queue
	runner        Runner
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

			if err := s.runner.Run(context.Background(), job); err != nil {
				errChan <- err
				continue
			}
		}
	}
}
