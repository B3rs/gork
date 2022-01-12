package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	columns = "id, status, queue, arguments, result, last_error, retry_count, options, created_at, updated_at, scheduled_at"
)

var (
	// ErrJobNotFound is returned when no job is found
	ErrJobNotFound = errors.New("job not found")
)

type Tx interface {
	List(ctx context.Context, limit int, offset int) ([]*Job, error)
	GetAndLock(ctx context.Context, queueName string) (*Job, error)
	Get(ctx context.Context, id string) (*Job, error)
	Update(ctx context.Context, job *Job) error
	Create(ctx context.Context, job *Job) error
	Deschedule(ctx context.Context, id string) error
	ScheduleImmediately(ctx context.Context, id string) error
	Commit() error
}

func NewTx(tx *sql.Tx) *Transaction {
	return &Transaction{Tx: tx}
}

type Transaction struct {
	*sql.Tx
}

func (tx *Transaction) GetAndLock(ctx context.Context, queueName string) (*Job, error) {
	lastError := &sql.NullString{}
	job := &Job{}

	err := tx.QueryRowContext(ctx, `SELECT `+columns+`
	FROM jobs 
	WHERE status = $1
		AND scheduled_at <= $2
		AND queue = $3
	ORDER BY scheduled_at ASC 
	FOR UPDATE SKIP LOCKED
	LIMIT 1 `, StatusScheduled, time.Now(), queueName).
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
		return nil, ErrNoJobsAvailable
	}
	if err != nil {
		return nil, err
	}

	job.LastError = lastError.String

	return job, nil
}

// Update the job in the database
func (tx *Transaction) Update(ctx context.Context, job *Job) error {
	res := sql.NullString{String: string(job.Result), Valid: true}
	if len(job.Result) == 0 {
		res.Valid = false
	}

	_, err := tx.ExecContext(ctx, `UPDATE jobs
	SET
		status=$1, 
		result=$2, 
		last_error=$3, 
		retry_count=$4,
		scheduled_at=$5,
		updated_at=now()
	WHERE id = $6`, job.Status, res, job.LastError, job.RetryCount, job.ScheduledAt, job.ID)
	return err
}

// Create the job in the database
func (tx *Transaction) Create(ctx context.Context, job *Job) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO jobs (id, queue, status, arguments, options, scheduled_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		job.ID, job.Queue, job.Status, job.Arguments, job.Options, job.ScheduledAt)
	return err
}

// Deschedule the job
func (tx *Transaction) Deschedule(ctx context.Context, id string) error {

	_, err := tx.ExecContext(ctx, `UPDATE jobs 
	SET 
		updated_at=now(), 
		scheduled_at=NULL, 
		status=$1 
	WHERE 
		id = $2 AND 
		status = $3`, StatusCanceled, id, StatusScheduled)
	return err
}

// ScheduleImmediately the job
func (tx *Transaction) ScheduleImmediately(ctx context.Context, id string) error {
	_, err := tx.ExecContext(ctx, `UPDATE jobs 
	SET 
		updated_at=now(), 
		scheduled_at=now(), 
		status=$1 
	WHERE 
		id = $2`, StatusScheduled, id)
	return err
}

func (tx *Transaction) List(ctx context.Context, limit, offset int) ([]*Job, error) {
	jobs := []*Job{}

	rows, err := tx.QueryContext(ctx, `SELECT `+columns+`
	FROM jobs 
	ORDER BY scheduled_at ASC 
	LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		lastError := &sql.NullString{}
		job := &Job{}
		rows.Scan(
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
		job.LastError = lastError.String
		fmt.Println("job", job)
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (tx *Transaction) Get(ctx context.Context, id string) (*Job, error) {
	lastError := &sql.NullString{}
	job := &Job{}
	const getSQL = `SELECT ` + columns + `
		FROM jobs 
		WHERE id = $1`

	err := tx.QueryRowContext(ctx, getSQL, id).
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
