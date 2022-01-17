package workers

import (
	"context"
	"time"

	"github.com/B3rs/gork/jobs"
)

func newDequeuer(q Queue, w Worker, sleepInterval time.Duration) *dequeuer {
	return &dequeuer{
		queue:         q,
		handler:       newHandler(w, q),
		sleepInterval: sleepInterval,
	}
}

type dequeuer struct {
	queue         Queue
	handler       Handler
	sleepInterval time.Duration
}

// Run starts a loop that polls the queue for jobs and passes them to the job handler
func (d *dequeuer) Run(ctx context.Context, errChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			job, err := d.queue.Dequeue(ctx)
			if err == jobs.ErrNoJobsAvailable {
				time.Sleep(d.sleepInterval)
				continue
			}
			if err != nil {
				errChan <- err
				continue
			}

			if err := d.handler.Handle(context.Background(), job); err != nil {
				errChan <- err
				continue
			}
		}
	}
}
