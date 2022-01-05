package workers

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

var (
	// ErrJobNotFound is returned when no job is found
	ErrJobNotFound = errors.New("job not found")
	// ErrQueueIsClosed is returned when the queue is closing
	ErrQueueIsClosed = errors.New("queue is closing")
)

func NewQueue(name string, db *sql.DB) *Queue {
	return &Queue{
		name: name,
		db:   db,
	}
}

type Queue struct {
	sync.Mutex

	db   *sql.DB
	name string
}

const acquireSQL = `SELECT id, status, queue, arguments, result, last_error, created_at, updated_at, scheduled_at
FROM jobs 
WHERE status = $1
	AND scheduled_at <= now()
	AND queue = $2
ORDER BY scheduled_at ASC 
FOR UPDATE SKIP LOCKED
LIMIT 1 `

// AcquireJob acquires a job from the database and locks it until the transaction is committed or rolled back
func (q *Queue) AcquireJobs(ctx context.Context, count int) ([]jobs.Job, error) {

	tx, err := q.db.Begin()
	if err != nil {
		return nil, err
	}

	job := jobs.Job{Tx: tx}

	lastError := &sql.NullString{}

	// TODO implement real scan on query with variable limit
	err = tx.QueryRowContext(ctx, acquireSQL, jobs.StatusScheduled, q.name).
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
		return nil, tx.Commit()
	}
	if err != nil {
		return nil, err
	}

	job.LastError = lastError.String

	return []jobs.Job{job}, err
}
