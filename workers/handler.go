package workers

import (
	"context"

	"github.com/B3rs/gork/jobs"
)

type updater interface {
	Update(context.Context, jobs.Job) error
}

func newHandler(worker Worker, updater updater) *handler {
	return &handler{
		worker:  worker,
		updater: updater,
	}
}

// handler runs a job in a worker, managing it's execution, errors and results.
type handler struct {
	worker  Worker
	updater updater
}

// Handle a job execution
func (h *handler) Handle(ctx context.Context, job jobs.Job) error {

	res, err := h.worker.Execute(ctx, job)
	if err != nil {
		return h.fail(ctx, job, err)
	}

	return h.success(ctx, job, res)
}

func (h *handler) fail(ctx context.Context, job jobs.Job, err error) error {

	job = job.SetLastError(err)

	if job.ShouldRetry() {
		retryAt := now().Add(job.Options.RetryInterval)
		job = job.ScheduleRetry(retryAt)
		return h.updater.Update(ctx, job)
	}

	job = job.SetStatus(jobs.StatusFailed)

	if err := h.updater.Update(ctx, job); err != nil {
		return err
	}

	return execOnFailureCallback(ctx, h.worker, job)
}

func (h *handler) success(ctx context.Context, job jobs.Job, res interface{}) error {
	job = job.SetStatus(jobs.StatusCompleted)

	var err error
	if job, err = job.SetResult(res); err != nil {
		job = job.SetLastError(err)
	}

	return h.updater.Update(ctx, job)
}

type failer interface {
	OnFailure(context.Context, jobs.Job) error
}

func execOnFailureCallback(ctx context.Context, worker Worker, job jobs.Job) error {
	f, ok := worker.(failer)
	if !ok {
		return nil
	}

	return f.OnFailure(ctx, job)
}
