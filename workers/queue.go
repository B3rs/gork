package workers

import (
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

	closing bool
}

// Close closes the queue
func (q *Queue) Close() {
	q.Lock()
	q.closing = true
	q.Unlock()
}

// isClosed returns true if the queue is closing
func (q *Queue) isClosed() bool {
	q.Lock()
	defer q.Unlock()
	return q.closing
}

const acquireSQL = `SELECT id, status, queue, arguments, result, last_error, created_at, updated_at, scheduled_at
FROM jobs 
WHERE status = 'scheduled' 
	AND scheduled_at <= now()
	AND queue = $1
ORDER BY scheduled_at ASC 
FOR UPDATE SKIP LOCKED
LIMIT 1 `

// AcquireJob acquires a job from the database and locks it until the transaction is committed or rolled back
func (q *Queue) AcquireJob() (jobs.Job, error) {
	if q.isClosed() {
		return jobs.Job{}, ErrQueueIsClosed
	}

	tx, err := q.db.Begin()
	if err != nil {
		return jobs.Job{}, err
	}

	job := jobs.Job{Tx: tx}

	lastError := &sql.NullString{}

	err = tx.QueryRow(acquireSQL, q.name).
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
		if err := tx.Commit(); err != nil {
			return jobs.Job{}, err
		}
		return jobs.Job{}, ErrJobNotFound
	}

	job.LastError = lastError.String

	return job, err
}
