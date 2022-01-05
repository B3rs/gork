package workers

import (
	"context"
	"database/sql"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

type acquirer interface {
	AcquireJobs(ctx context.Context, limit int) ([]jobs.Job, error)
}

type poller struct {
	acquirer acquirer
	sleep    time.Duration

	jobsChan chan jobs.Job
}

func newPoller(db *sql.DB, queueName string, sleepInterval time.Duration) poller {
	return poller{
		acquirer: NewQueue(queueName, db),
		sleep:    sleepInterval,
		jobsChan: make(chan jobs.Job),
	}
}

func (p poller) jobs() <-chan jobs.Job {
	return p.jobsChan
}

func (p poller) stop() {
	close(p.jobsChan)
}

func (p *poller) start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Acquire job and lock it
			jobs, err := p.acquirer.AcquireJobs(ctx, 1)
			if err != nil {
				return err
			}
			if len(jobs) == 0 {
				time.Sleep(p.sleep)
				continue
			}

			for _, job := range jobs {
				p.jobsChan <- job
			}
		}
	}

}
