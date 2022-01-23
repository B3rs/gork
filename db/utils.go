package db

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/jobs"
)

func queryJob(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (jobs.Job, error) {
	job := &job{}

	err := tx.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(job.ScanDestinations()...)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return jobs.Job{}, jobs.ErrJobNotFound
	default:
		return jobs.Job{}, err
	}

	return job.ToJob(), nil
}

func queryJobs(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) ([]jobs.Job, error) {
	jobs := []jobs.Job{}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		job := &job{}
		rows.Scan(job.ScanDestinations()...)
		jobs = append(jobs, job.ToJob())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func exec(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) error {
	_, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)
	return err
}
