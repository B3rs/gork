package client

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/db"
	"github.com/B3rs/gork/jobs"
)

func NewClient(database *sql.DB) *Client {
	return &Client{
		db: database,
		txClientFactory: func(tx *sql.Tx) TxClient {
			return newTx(tx)
		},
		txWrapper: db.NewTxWrapper(database),
	}
}

type txWrapper interface {
	WrapTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) (interface{}, error)) (interface{}, error)
}

type Client struct {
	db              *sql.DB
	txClientFactory func(*sql.Tx) TxClient
	txWrapper       txWrapper
}

// WithTx returns a new client with the given transaction
func (c Client) WithTx(tx *sql.Tx) TxClient {
	return c.txClientFactory(tx)
}

// Schedule schedules a job in the queue to be executed as soon as possible
func (c Client) Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error {
	_, err := c.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return nil, c.txClientFactory(tx).Schedule(ctx, id, queueName, arguments, options...)
	})
	return err
}

// Cancel cancels a job in the queue if not already executed
func (c Client) Cancel(ctx context.Context, id string) error {
	_, err := c.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return nil, c.txClientFactory(tx).Cancel(ctx, id)
	})
	return err
}

// ScheduleNow reschedules a job in the queue to be executed immediately
func (c Client) ScheduleNow(ctx context.Context, id string) error {
	_, err := c.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return nil, c.txClientFactory(tx).ScheduleNow(ctx, id)
	})
	return err
}

// Get returns job with the given id
func (c Client) Get(ctx context.Context, id string) (jobs.Job, error) {
	res, err := c.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return c.txClientFactory(tx).Get(ctx, id)
	})
	if err != nil {
		return jobs.Job{}, err
	}
	return res.(jobs.Job), nil
}
