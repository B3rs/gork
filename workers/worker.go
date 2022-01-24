package workers

//go:generate mockgen -destination=./mocks_test.go -package=workers github.com/B3rs/gork/workers Queue,Worker,Handler,Spawner

import (
	"context"
	"time"

	"github.com/B3rs/gork/jobs"
)

var (
	now = time.Now
)

// Queue is a queue of jobs.
type Queue interface {
	Dequeue(ctx context.Context) (jobs.Job, error)
	RequeueTimedOutJobs(ctx context.Context, timeout time.Duration) error
}

// Handler handles job execution, errors and results.
type Handler interface {
	Handle(ctx context.Context, job jobs.Job) error
}

type Spawner interface {
	Spawn(runner)
	Wait()
	Shutdown()
	Done() <-chan struct{}
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
