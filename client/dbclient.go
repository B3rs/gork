package client

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/jobs"
)

func NewDBClient(db *sql.DB) *DBClient {
	return &DBClient{
		db: db,
		txClientFactory: func(tx *sql.Tx) Client {
			return NewTxClient(tx)
		},
	}
}

type DBClient struct {
	db              *sql.DB
	txClientFactory func(*sql.Tx) Client
}

// Schedule schedules a job in the queue to be executed as soon as possible
func (c DBClient) Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := c.txClientFactory(tx).Schedule(ctx, id, queueName, arguments, options...); err != nil {
		return err
	}
	return tx.Commit()
}

// Cancel cancels a job in the queue if not already executed
func (c DBClient) Cancel(ctx context.Context, id string) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := c.txClientFactory(tx).Cancel(ctx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// ForceRetry reschedules a job in the queue to be executed immediately
func (c DBClient) ForceRetry(ctx context.Context, id string) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := c.txClientFactory(tx).ForceRetry(ctx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// GetAll returns jobs starting from the given offset
func (c DBClient) GetAll(ctx context.Context, page, limit int) ([]*jobs.Job, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	jobs, err := c.txClientFactory(tx).GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return jobs, tx.Commit()
}

// Get returns job with the given id
func (c DBClient) Get(ctx context.Context, id string) (*jobs.Job, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	jobs, err := c.txClientFactory(tx).Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return jobs, tx.Commit()
}
