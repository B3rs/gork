package workers

import (
	"context"

	"github.com/B3rs/gork/jobs"
)

func newRunner(worker Worker, updater Queue) *runner {
	return &runner{
		worker:  worker,
		updater: updater,
	}
}

// runner runs a job in a worker, managing it's execution, errors and results.
type runner struct {
	worker  Worker
	updater Queue
}

func (r *runner) Run(ctx context.Context, job *jobs.Job) error {

	// explicity copy the job to avoid user provided function from modifying a job
	res, err := r.worker.Execute(ctx, *job)
	if err != nil {
		if job.ShouldRetry() {
			retryAt := now().Add(job.Options.RetryInterval)
			job.ScheduleRetry(retryAt)
		} else {
			job.SetStatus(jobs.StatusFailed)
		}
		job.SetLastError(err)

		return r.updater.Update(ctx, job)
	}

	job.SetStatus(jobs.StatusCompleted)
	if err := job.SetResult(res); err != nil {
		job.SetLastError(err)
	}

	return r.updater.Update(ctx, job)
}
