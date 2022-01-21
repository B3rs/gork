package db

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/jobs"
)

func NewTx(tx *sql.Tx) *Transaction {
	return &Transaction{Tx: tx}
}

type Transaction struct {
	*sql.Tx
}

// Update the job in the database
func (tx *Transaction) Update(ctx context.Context, job *jobs.Job) error {

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
func (tx *Transaction) Create(ctx context.Context, job *jobs.Job) error {
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
		status=$1 
	WHERE 
		id = $2 AND 
		status = $3`, jobs.StatusCanceled, id, jobs.StatusScheduled)
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
		id = $2`, jobs.StatusScheduled, id)
	return err
}

func (tx *Transaction) List(ctx context.Context, limit, offset int, search string) ([]*jobs.Job, error) {
	jobs := []*jobs.Job{}

	columns := jobRow{}.StringColumns()
	rows, err := tx.QueryContext(ctx, `SELECT `+columns+`
	FROM jobs
	WHERE id LIKE '%' || $1 || '%' 
	ORDER BY scheduled_at DESC 
	LIMIT $2 OFFSET $3`, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		job := &jobRow{}
		rows.Scan(job.ScanDestinations()...)

		jobs = append(jobs, job.ToJob())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (tx *Transaction) Get(ctx context.Context, id string) (*jobs.Job, error) {

	columns := jobRow{}.StringColumns()
	getSQL := `SELECT ` + columns + `
		FROM jobs 
		WHERE id = $1`

	job := &jobRow{}
	err := tx.QueryRowContext(ctx, getSQL, id).
		Scan(job.ScanDestinations()...)

	if err == sql.ErrNoRows {
		return nil, jobs.ErrJobNotFound
	}
	if err != nil {
		return nil, err
	}

	return job.ToJob(), nil
}
