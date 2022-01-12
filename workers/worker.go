package workers

import (
	"context"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

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
