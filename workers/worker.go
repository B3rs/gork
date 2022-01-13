package workers

//go:generate mockgen -destination=./mocks_test.go -package=workers github.com/B3rs/gork/workers Requeuer,Queue,Runner,Worker,jobUpdater

import (
	"context"
	"time"

	"github.com/B3rs/gork/jobs"
)

var (
	now = time.Now
)

// Requeuer re-queues timed out jobs.
type Requeuer interface {
	RequeueTimedOutJobs(ctx context.Context, timeout time.Duration) error
}

// Queue is a queue of jobs.
type Queue interface {
	Dequeue(ctx context.Context) (*jobs.Job, error)
	Update(ctx context.Context, job *jobs.Job) error
}

// Runner runs a job in a worker, managing it's execution, errors and results.
type Runner interface {
	Run(ctx context.Context, job *jobs.Job) error
}

// WorkerFunc is a function that can be used as a worker.
type WorkerFunc func(ctx context.Context, job jobs.Job) (interface{}, error)

// Worker is the interface that must be implemented by workers.
type Worker interface {
	Execute(ctx context.Context, job jobs.Job) (interface{}, error)
}

// funcWorker is just a wrapper for a WorkerFunc.
type funcWorker struct {
	f WorkerFunc
}

// Execute runs the worker function for the given job.
func (f funcWorker) Execute(ctx context.Context, job jobs.Job) (interface{}, error) {
	return f.f(ctx, job)
}
