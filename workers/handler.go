package workers

import (
	"context"

	"github.com/B3rs/gork/jobs"
)

func newHandler(worker Worker, updater Queue) *handler {
	return &handler{
		worker:  worker,
		updater: updater,
	}
}

// handler runs a job in a worker, managing it's execution, errors and results.
type handler struct {
	worker  Worker
	updater Queue
}

// Handle a job execution
func (h *handler) Handle(ctx context.Context, job *jobs.Job) error {

	// explicity copy the job to avoid user provided function from modifying a job
	res, err := h.worker.Execute(ctx, *job)
	if err != nil {
		if job.ShouldRetry() {
			retryAt := now().Add(job.Options.RetryInterval)
			job.ScheduleRetry(retryAt)
		} else {
			job.SetStatus(jobs.StatusFailed)
		}
		job.SetLastError(err)

		return h.updater.Update(ctx, job)
	}

	job.SetStatus(jobs.StatusCompleted)
	if err := job.SetResult(res); err != nil {
		job.SetLastError(err)
	}

	return h.updater.Update(ctx, job)
}
