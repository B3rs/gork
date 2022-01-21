package client

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/db"
	"github.com/B3rs/gork/jobs"
)

type Client interface {
	Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error
	Cancel(ctx context.Context, id string) error
	ForceRetry(ctx context.Context, id string) error
	GetAll(ctx context.Context, page, limit int, search string) ([]*jobs.Job, error)
	Get(ctx context.Context, id string) (*jobs.Job, error)
}

func NewTx(tx *sql.Tx) Client {
	return Tx{tx: db.NewTx(tx)}
}

type Tx struct {
	tx db.Tx
}

// Schedule schedules a job in the queue to be executed as soon as possible
func (c Tx) Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error {
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

	return c.tx.Create(ctx, job)
}

// Cancel cancels a job in the queue if not already executed
func (c Tx) Cancel(ctx context.Context, id string) error {
	return c.tx.Deschedule(ctx, id)
}

// ForceRetry reschedules a job in the queue to be executed immediately
func (c Tx) ForceRetry(ctx context.Context, id string) error {
	return c.tx.ScheduleImmediately(ctx, id)
}

// GetAll returns jobs starting from the given offset
func (c Tx) GetAll(ctx context.Context, page, limit int, search string) ([]*jobs.Job, error) {
	return c.tx.List(ctx, limit, page*limit, search)
}

// Get returns job with the given id
func (c Tx) Get(ctx context.Context, id string) (*jobs.Job, error) {
	return c.tx.Get(ctx, id)
}
