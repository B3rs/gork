package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const (
	// StatusCompleted is the status of a job that has been completed
	StatusCompleted = "completed"
	// StatusScheduled is the status of a job that has been scheduled
	StatusScheduled = "scheduled"
	// StatusFailed is the status of a job that has failed
	StatusFailed = "failed"
	// StatusCanceled is the status of a job that has been canceled
	StatusCanceled = "canceled"
)

// Job could use generics for params and result
type Job struct {
	ID          string
	Queue       string
	Status      string
	Arguments   []byte
	Result      []byte
	LastError   string
	ScheduledAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Tx *sql.Tx
}

func (j Job) ParseArguments(dest interface{}) error {
	return json.Unmarshal(j.Arguments, dest)
}

// SetStatus sets the status of a job
func (j *Job) SetStatus(ctx context.Context, status string) error {
	_, err := j.Tx.ExecContext(ctx,
		`UPDATE jobs
SET status = $2,
updated_at = now()
WHERE id = $1`,
		j.ID,
		status,
	)

	return err
}

// SetResult sets the result of a job
func (j *Job) SetResult(ctx context.Context, result interface{}) error {
	if result == nil {
		return nil
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		return err
	}

	_, err = j.Tx.ExecContext(ctx,
		`UPDATE jobs
SET result = $2,
updated_at = now()
WHERE id = $1`,
		j.ID,
		encoded,
	)
	return err
}

// SetLastError sets the last error of a job
func (j *Job) SetLastError(ctx context.Context, e error) error {
	_, err := j.Tx.ExecContext(ctx,
		`UPDATE jobs
SET last_error = $2,
updated_at = now()
WHERE id = $1`,
		j.ID,
		e.Error(),
	)
	return err
}

// Commit commits the transaction
func (j *Job) Commit() error {
	return j.Tx.Commit()
}

// Rollback rolls back the transaction
func (j *Job) Rollback() error {
	return j.Tx.Rollback()
}
