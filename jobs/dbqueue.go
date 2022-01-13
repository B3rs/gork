package jobs

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	// ErrNoJobsAvailable is returned when no jobs are available
	ErrNoJobsAvailable = errors.New("no jobs available")
)

func NewQueue(db *sql.DB, name string) *Queue {
	return &Queue{
		db:   db,
		name: name,
	}
}

type Queue struct {
	db   *sql.DB
	name string
}

func (q *Queue) Dequeue(ctx context.Context) (*Job, error) {
	job := &jobRow{}

	tx, err := q.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() { _ = tx.Rollback() }()

	query := `UPDATE jobs
	SET 
		status=$1, 
		updated_at=now(),
		started_at=now()
	WHERE 
		id = (
			SELECT id
			FROM jobs 
			WHERE status = $2
				AND scheduled_at <= now()
				AND queue = $3
			ORDER BY scheduled_at ASC 
			FOR UPDATE SKIP LOCKED
			LIMIT 1 
		)
	RETURNING ` + job.stringColumns()

	err = tx.QueryRowContext(
		ctx,
		query,
		StatusInitialized,
		StatusScheduled,
		q.name,
	).Scan(job.scanDestinations()...)

	if err == sql.ErrNoRows {
		if err := tx.Commit(); err != nil {
			return nil, err
		}

		return nil, ErrNoJobsAvailable
	}
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return job.toJob(), nil
}

// Update the job in the database
func (q *Queue) Update(ctx context.Context, job *Job) error {

	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	defer func() { _ = tx.Commit() }()

	res := sql.NullString{String: string(job.Result), Valid: true}
	if len(job.Result) == 0 {
		res.Valid = false
	}

	_, err = tx.ExecContext(ctx, `UPDATE jobs
	SET
		status=$1, 
		result=$2, 
		last_error=$3, 
		retry_count=$4,
		scheduled_at=$5,
		updated_at=now()
	WHERE id = $6`, job.Status, res, job.LastError, job.RetryCount, job.ScheduledAt, job.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (q *Queue) RequeueTimedOutJobs(ctx context.Context, timeout time.Duration) error {

	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()

	query := `UPDATE jobs
	SET 
		status=$1, 
		started_at=null,
		retry_count=retry_count+1,
		updated_at=now()
	WHERE 
		started_at < $2 AND
		status = $3 AND
		queue = $4`

	_, err = tx.ExecContext(
		ctx,
		query,
		StatusScheduled,
		time.Now().Add(-timeout),
		StatusInitialized,
		q.name,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}