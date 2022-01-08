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

type Client interface {
	Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...OptionFunc) error
	Cancel(ctx context.Context, id string) error
	ForceRetry(ctx context.Context, id string) error
	GetAll(ctx context.Context, page, limit int) ([]*jobs.Job, error)
	Get(ctx context.Context, id string) (*jobs.Job, error)
}

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
