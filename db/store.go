package db

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/jobs"
)

func NewStore(db *sql.DB) *Store {
	return &Store{
		txWrapper:      NewTxWrapper(db),
		txStoreFactory: NewTx,
	}
}

//go:generate mockgen -destination=./store_mock.go -package=db -source=store.go
type JobsStore interface {
	Update(ctx context.Context, job jobs.Job) error
	Create(ctx context.Context, job jobs.Job) error
	Deschedule(ctx context.Context, id string) error
	ScheduleNow(ctx context.Context, id string) error
	Search(ctx context.Context, limit, offset int, search string) ([]jobs.Job, error)
	Get(ctx context.Context, id string) (jobs.Job, error)
	GetStatistics(ctx context.Context) (Statistics, error)
}

type Store struct {
	txWrapper      TxWrapper
	txStoreFactory func(tx *sql.Tx) TxStore
}

// Update the job in the database
func (s *Store) Update(ctx context.Context, job jobs.Job) error {
	_, err := s.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return nil, s.txStoreFactory(tx).Update(ctx, job)
	})
	return err
}

// Create the job in the database
func (s *Store) Create(ctx context.Context, job jobs.Job) error {
	_, err := s.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return nil, s.txStoreFactory(tx).Create(ctx, job)
	})
	return err
}

// Deschedule the job
func (s *Store) Deschedule(ctx context.Context, id string) error {
	_, err := s.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return nil, s.txStoreFactory(tx).Deschedule(ctx, id)
	})
	return err
}

// ScheduleNow the job
func (s *Store) ScheduleNow(ctx context.Context, id string) error {
	_, err := s.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return nil, s.txStoreFactory(tx).ScheduleNow(ctx, id)
	})
	return err
}

func (s *Store) Search(ctx context.Context, limit, offset int, search string) ([]jobs.Job, error) {
	res, err := s.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return s.txStoreFactory(tx).Search(ctx, limit, offset, search)
	})
	if err != nil {
		return nil, err
	}
	return res.([]jobs.Job), err
}

func (s *Store) Get(ctx context.Context, id string) (jobs.Job, error) {
	res, err := s.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		return s.txStoreFactory(tx).Get(ctx, id)
	})
	if err != nil {
		return jobs.Job{}, err
	}
	return res.(jobs.Job), err
}

func (s *Store) GetStatistics(ctx context.Context) (Statistics, error) {
	res, err := s.txWrapper.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		query := `SELECT *
		FROM (
			SELECT queue FROM jobs GROUP BY queue
		) AS queues
		JOIN LATERAL (
			SELECT count(*) AS scheduled FROM jobs WHERE queue = queues.queue AND status=$1
		) AS s ON true
		JOIN LATERAL (
			SELECT count(*) AS initialized FROM jobs WHERE queue = queues.queue AND status=$2
		) AS i ON true
		JOIN LATERAL (
			SELECT count(*) AS failed FROM jobs WHERE queue = queues.queue AND status=$3
		) AS f ON true
		JOIN LATERAL (
			SELECT count(*) AS completed FROM jobs WHERE queue = queues.queue AND status=$4
		) AS c ON true
		`
		qs, err := queryQueueStatistics(ctx, tx, query, jobs.StatusScheduled, jobs.StatusInitialized, jobs.StatusFailed, jobs.StatusCompleted)
		if err != nil {
			return nil, err
		}
		return Statistics{Queues: qs}, nil
	})
	if err != nil {
		return Statistics{}, err
	}
	return res.(Statistics), nil
}
