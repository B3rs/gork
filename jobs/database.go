package jobs

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	// ErrJobNotFound is returned when no job is found
	ErrJobNotFound = errors.New("job not found")
)

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

type Store struct {
	db *sql.DB
}

// GetAndLock gets a job from the database and locks it until the transaction is committed or rolled back
func (s *Store) Begin() (*Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

type Tx struct {
	*sql.Tx
}

const acquireSQL = `SELECT id, status, queue, arguments, result, last_error, retry_count, options, created_at, updated_at, scheduled_at
FROM jobs 
WHERE status = $1
	AND scheduled_at <= $2
	AND queue = $3
ORDER BY scheduled_at ASC 
FOR UPDATE SKIP LOCKED
LIMIT 1 `

func (tx *Tx) GetAndLock(ctx context.Context, queueName string) (*Job, error) {
	lastError := &sql.NullString{}
	job := &Job{}
	err := tx.QueryRowContext(ctx, acquireSQL, StatusScheduled, time.Now(), queueName).
		Scan(
			&job.ID,
			&job.Status,
			&job.Queue,
			&job.Arguments,
			&job.Result,
			lastError,
			&job.RetryCount,
			&job.Options,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.ScheduledAt,
		)

	if err == sql.ErrNoRows {
		return nil, ErrJobNotFound
	}
	if err != nil {
		return nil, err
	}

	job.LastError = lastError.String

	return job, nil
}

const updateSQL = `UPDATE jobs
SET
	status=$1, 
	result=$2, 
	last_error=$3, 
	retry_count=$4,
	scheduled_at=$5,
	updated_at=now()
WHERE id = $6`

// Update the job in the database
func (tx *Tx) Update(ctx context.Context, job *Job) error {
	res := sql.NullString{String: string(job.Result), Valid: true}
	if len(job.Result) == 0 {
		res.Valid = false
	}
	_, err := tx.ExecContext(ctx, updateSQL, job.Status, res, job.LastError, job.RetryCount, job.ScheduledAt, job.ID)
	return err
}

const createSQL = `INSERT INTO jobs 
	(id, queue, status, arguments, options, scheduled_at) 
VALUES 
	($1, $2, $3, $4, $5, $6)`

// Create the job in the database
func (tx *Tx) Create(ctx context.Context, job *Job) error {
	_, err := tx.ExecContext(ctx, createSQL, job.ID, job.Queue, job.Status, job.Arguments, job.Options, job.ScheduledAt)
	return err
}

const descheduleSQL = `UPDATE jobs 
SET 
	updated_at=now(), 
	scheduled_at=NULL, 
	status=$1 
WHERE 
	id = $2 AND 
	status = $3`

// Deschedule the job
func (tx *Tx) Deschedule(ctx context.Context, id string) error {
	_, err := tx.ExecContext(ctx, descheduleSQL, StatusCanceled, id, StatusScheduled)
	return err
}

const scheduleImmediatelySQL = `UPDATE jobs 
SET 
	updated_at=now(), 
	scheduled_at=now(), 
	status=$1 
WHERE 
	id = $2`

// ScheduleImmediately the job
func (tx *Tx) ScheduleImmediately(ctx context.Context, id string) error {
	_, err := tx.ExecContext(ctx, scheduleImmediatelySQL, StatusScheduled, id)
	return err
}
