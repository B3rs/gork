package db

import (
	"context"

	"github.com/B3rs/gork/jobs"
)

type Tx interface {
	Search(ctx context.Context, limit int, offset int, search string) ([]jobs.Job, error)
	Get(ctx context.Context, id string) (jobs.Job, error)
	Update(ctx context.Context, job jobs.Job) error
	Create(ctx context.Context, job jobs.Job) error
	Deschedule(ctx context.Context, id string) error
	ScheduleNow(ctx context.Context, id string) error
	Commit() error
}
