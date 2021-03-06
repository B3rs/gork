package db

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/jobs"
)

type TxStore interface {
	Search(ctx context.Context, limit int, offset int, search string) ([]jobs.Job, error)
	Get(ctx context.Context, id string) (jobs.Job, error)
	Update(ctx context.Context, job jobs.Job) error
	Create(ctx context.Context, job jobs.Job) error
	Deschedule(ctx context.Context, id string) error
	ScheduleNow(ctx context.Context, id string) error
	Commit() error
}

func NewTx(tx *sql.Tx) TxStore {
	return &Tx{Tx: tx}
}

type Tx struct {
	*sql.Tx
}

// Update the job in the database
func (tx *Tx) Update(ctx context.Context, job jobs.Job) error {
	return exec(ctx, tx.Tx, `UPDATE jobs
	SET
		status=$1, 
		result=$2, 
		last_error=$3, 
		retry_count=$4,
		scheduled_at=$5,
		updated_at=now()
	WHERE id = $6`,
		job.Status,
		job.Result,
		job.LastError,
		job.RetryCount,
		job.ScheduledAt,
		job.ID,
	)
}

// Create the job in the database
func (tx *Tx) Create(ctx context.Context, job jobs.Job) error {
	return exec(ctx, tx.Tx, `INSERT INTO jobs (id, queue, status, arguments, options, scheduled_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		job.ID, job.Queue, job.Status, job.Arguments, job.Options, job.ScheduledAt)
}

// Deschedule the job
func (tx *Tx) Deschedule(ctx context.Context, id string) error {
	return exec(ctx, tx.Tx, `UPDATE jobs 
	SET 
		updated_at=now(), 
		status=$1 
	WHERE 
		id = $2 AND 
		status = $3`, jobs.StatusCanceled, id, jobs.StatusScheduled)
}

// ScheduleNow the job
func (tx *Tx) ScheduleNow(ctx context.Context, id string) error {
	return exec(ctx, tx.Tx, `UPDATE jobs 
	SET 
		updated_at=now(), 
		scheduled_at=now(), 
		status=$1 
	WHERE 
		id = $2`, jobs.StatusScheduled, id)
}

func (tx *Tx) Search(ctx context.Context, limit, offset int, search string) ([]jobs.Job, error) {
	if search != "" {
		return queryJobs(ctx, tx.Tx, `SELECT `+jobStringColumns+`
		FROM jobs
		WHERE id LIKE '%' || $1 || '%' 
		ORDER BY scheduled_at DESC 
		LIMIT $2 OFFSET $3`, search, limit, offset)
	}
	return queryJobs(ctx, tx.Tx, `SELECT `+jobStringColumns+`
	FROM jobs
	ORDER BY scheduled_at DESC 
	LIMIT $1 OFFSET $2`, limit, offset)
}

func (tx *Tx) Get(ctx context.Context, id string) (jobs.Job, error) {
	return queryJob(ctx, tx.Tx, `SELECT `+jobStringColumns+`
	FROM jobs 
	WHERE id = $1`, id)
}
