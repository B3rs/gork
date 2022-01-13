package client

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/jobs"
)

type Client interface {
	Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error
	Cancel(ctx context.Context, id string) error
	ForceRetry(ctx context.Context, id string) error
	GetAll(ctx context.Context, page, limit int) ([]*jobs.Job, error)
	Get(ctx context.Context, id string) (*jobs.Job, error)
}

func NewTxClient(tx *sql.Tx) Client {
	return TxClient{tx: jobs.NewTx(tx)}
}

type TxClient struct {
	tx jobs.Tx
}

// Schedule schedules a job in the queue to be executed as soon as possible
func (c TxClient) Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error {
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
func (c TxClient) Cancel(ctx context.Context, id string) error {
	return c.tx.Deschedule(ctx, id)
}

// ForceRetry reschedules a job in the queue to be executed immediately
func (c TxClient) ForceRetry(ctx context.Context, id string) error {
	return c.tx.ScheduleImmediately(ctx, id)
}

// GetAll returns jobs starting from the given offset
func (c TxClient) GetAll(ctx context.Context, page, limit int) ([]*jobs.Job, error) {
	return c.tx.List(ctx, limit, page*limit)
}

// Get returns job with the given id
func (c TxClient) Get(ctx context.Context, id string) (*jobs.Job, error) {
	return c.tx.Get(ctx, id)
}
