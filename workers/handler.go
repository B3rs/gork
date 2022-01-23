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
func (h *handler) Handle(ctx context.Context, job jobs.Job) error {

	res, err := h.worker.Execute(ctx, job)
	if err != nil {
		if job.ShouldRetry() {
			retryAt := now().Add(job.Options.RetryInterval)
			job = job.ScheduleRetry(retryAt)
		} else {
			job = job.SetStatus(jobs.StatusFailed)
		}

		job = job.SetLastError(err)
		return h.updater.Update(ctx, job)
	}

	job = job.SetStatus(jobs.StatusCompleted)

	if job, err = job.SetResult(res); err != nil {
		job = job.SetLastError(err)
	}

	return h.updater.Update(ctx, job)
}
