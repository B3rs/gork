package client

import (
	"context"
	"database/sql"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

var (
	defaultOptions = []OptionFunc{
		WithRetryInterval(10 * time.Second),
		scheduleImmediately(),
	}
)

// Schedule schedules a job in the queue to be executed as soon as possible
func Schedule(ctx context.Context, tx *sql.Tx, id string, queueName string, arguments interface{}, options ...OptionFunc) error {
	job := &jobs.Job{
		ID:     id,
		Queue:  queueName,
		Status: jobs.StatusScheduled,
	}

	if err := job.SetArguments(arguments); err != nil {
		return err
	}

	for _, opt := range append(defaultOptions, options...) {
		job = opt(job)
	}

	return (&jobs.Tx{Tx: tx}).Create(ctx, job)
}

// Cancel cancels a job in the queue if not already executed
func Cancel(ctx context.Context, tx *sql.Tx, id string) error {
	return (&jobs.Tx{Tx: tx}).Deschedule(ctx, id)
}

// ForceRetry reschedules a job in the queue to be executed immediately
func ForceRetry(ctx context.Context, tx *sql.Tx, id string) error {
	return (&jobs.Tx{Tx: tx}).ScheduleImmediately(ctx, id)
}

type OptionFunc func(j *jobs.Job) *jobs.Job

func WithMaxRetries(tries int) OptionFunc {
	return func(j *jobs.Job) *jobs.Job {
		j.Options.MaxRetries = tries
		return j
	}
}

func WithRetryInterval(interval time.Duration) OptionFunc {
	return func(j *jobs.Job) *jobs.Job {
		j.Options.RetryInterval = interval
		return j
	}
}

func WithScheduleTime(t time.Time) OptionFunc {
	return func(j *jobs.Job) *jobs.Job {
		j.ScheduledAt = t
		return j
	}
}

func scheduleImmediately() OptionFunc {
	return func(j *jobs.Job) *jobs.Job {
		j.ScheduledAt = time.Now()
		return j
	}
}
