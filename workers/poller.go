package workers

import (
	"context"
	"time"

	"github.com/B3rs/gork/jobs"
)

func newPoller(q Queue, updater updater, w Worker, sleepInterval time.Duration) *poller {
	return &poller{
		queue:         q,
		handler:       newHandler(w, updater),
		sleepInterval: sleepInterval,
	}
}

type poller struct {
	queue         Queue
	handler       Handler
	sleepInterval time.Duration
}

// Run starts a loop that polls the queue for jobs and passes them to the job handler
func (d *poller) Run(ctx context.Context, errChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			job, err := d.queue.Pop(ctx)
			if err == jobs.ErrJobNotFound {
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
