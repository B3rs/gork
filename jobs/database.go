package jobs

import (
	"context"
	"database/sql"
	"errors"
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

const acquireSQL = `SELECT id, status, queue, arguments, result, last_error, created_at, updated_at, scheduled_at
FROM jobs 
WHERE status = $1
	AND scheduled_at <= now()
	AND queue = $2
ORDER BY scheduled_at ASC 
FOR UPDATE SKIP LOCKED
LIMIT 1 `

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

func (tx *Tx) GetAndLock(ctx context.Context, queueName string) (*Job, error) {
	lastError := &sql.NullString{}
	job := &Job{}
	err := tx.QueryRowContext(ctx, acquireSQL, StatusScheduled, queueName).
		Scan(
			&job.ID,
			&job.Status,
			&job.Queue,
			&job.Arguments,
			&job.Result,
			lastError,
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
	updated_at=now()
WHERE id = $4`

// Update the job in the database
func (tx *Tx) Update(ctx context.Context, job *Job) error {
	res := sql.NullString{String: string(job.Result), Valid: true}
	if len(job.Result) == 0 {
		res.Valid = false
	}
	_, err := tx.ExecContext(ctx, updateSQL, job.Status, res, job.LastError, job.ID)
	return err
}
