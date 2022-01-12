package workers

import (
	"context"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

func newExecutor(
	worker Worker,
	update updateJobFunc,
) *executor {
	return &executor{
		worker: worker,
		update: update,
	}
}

// executor implements the logic necessary to execute a worker.
type executor struct {
	worker Worker
	update updateJobFunc
}

type updateJobFunc func(ctx context.Context, job *jobs.Job) error

func (e *executor) execute(ctx context.Context, job *jobs.Job) error {

	// explicity copy the job to avoid user provided function to modify the job
	res, err := e.worker.Execute(ctx, *job)
	if err != nil {
		if job.ShouldRetry() {
			retryAt := time.Now().Add(job.Options.RetryInterval)
			job.ScheduleRetry(retryAt)
		} else {
			job.SetStatus(jobs.StatusFailed)
		}
		job.SetLastError(err)

		return e.update(ctx, job)
	}

	job.SetStatus(jobs.StatusCompleted)
	if err := job.SetResult(res); err != nil {
		job.SetLastError(err)
	}

	return e.update(ctx, job)
}
