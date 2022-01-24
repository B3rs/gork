package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/B3rs/gork/jobs"
)

func NewQueue(db *sql.DB, name string) *Queue {
	return &Queue{
		db:        db,
		name:      name,
		TxWrapper: NewTxWrapper(db),
	}
}

type Queue struct {
	TxWrapper

	db   *sql.DB
	name string
}

func (q *Queue) Dequeue(ctx context.Context) (jobs.Job, error) {
	res, err := q.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
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
		RETURNING ` + jobStringColumns

		return queryJob(
			ctx,
			tx,
			query,
			jobs.StatusInitialized,
			jobs.StatusScheduled,
			q.name,
		)
	})
	if err != nil {
		return jobs.Job{}, err
	}
	return res.(jobs.Job), nil
}

// Update the job in the database
func (q *Queue) Update(ctx context.Context, job jobs.Job) error {
	_, err := q.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {

		return nil, exec(ctx, tx, `UPDATE jobs
		SET
			status=$1, 
			result=$2, 
			last_error=$3, 
			retry_count=$4,
			scheduled_at=$5,
			updated_at=now()
		WHERE id = $6`,
			job.Status,
			// sql.NullString{String: string(job.Result), Valid: len(job.Result) != 0},
			job.Result,
			job.LastError,
			job.RetryCount,
			job.ScheduledAt,
			job.ID,
		)
	})
	return err
}

func (q *Queue) RequeueTimedOutJobs(ctx context.Context, timeout time.Duration) error {
	_, err := q.WrapTx(ctx, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
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

		return nil, exec(
			ctx,
			tx,
			query, jobs.StatusScheduled,
			time.Now().Add(-timeout),
			jobs.StatusInitialized,
			q.name)
	})

	return err
}
