package client

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/db"
	"github.com/B3rs/gork/jobs"
)

//go:generate mockgen -destination=./txmocks_test.go -package=client -source=tx.go

type TxClient interface {
	Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error
	Cancel(ctx context.Context, id string) error
	ScheduleNow(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (jobs.Job, error)
}

func newTx(tx *sql.Tx) TxClient {
	return txClient{tx: db.NewTx(tx)}
}

type dbTx interface {
	Get(ctx context.Context, id string) (jobs.Job, error)
	Create(ctx context.Context, job jobs.Job) error
	Deschedule(ctx context.Context, id string) error
	ScheduleNow(ctx context.Context, id string) error
}

type txClient struct {
	tx dbTx
}

// Schedule schedules a job in the queue to be executed as soon as possible
func (c txClient) Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error {
	job := jobs.Job{
		ID:     id,
		Queue:  queueName,
		Status: jobs.StatusScheduled,
	}

	var err error
	if job, err = job.SetArguments(arguments); err != nil {
		return err
	}

	for _, opt := range append(defaultOptions, options...) {
		job = opt(job)
	}

	return c.tx.Create(ctx, job)
}

// Cancel cancels a job in the queue if not already executed
func (c txClient) Cancel(ctx context.Context, id string) error {
	return c.tx.Deschedule(ctx, id)
}

// ForceRetry reschedules a job in the queue to be executed immediately
func (c txClient) ScheduleNow(ctx context.Context, id string) error {
	return c.tx.ScheduleNow(ctx, id)
}

// Get returns job with the given id
func (c txClient) Get(ctx context.Context, id string) (jobs.Job, error) {
	return c.tx.Get(ctx, id)
}
