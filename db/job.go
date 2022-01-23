package db

import (
	"database/sql"
	"strings"
	"time"

	"github.com/B3rs/gork/jobs"
)

var (
	jobColumns = []string{
		"id",
		"status",
		"queue",
		"arguments",
		"result",
		"last_error",
		"retry_count",
		"options",
		"created_at",
		"updated_at",
		"scheduled_at",
	}

	jobStringColumns = strings.Join(jobColumns, ", ")
)

type job struct {
	ID          string         `json:"id"`
	Queue       string         `json:"queue"`
	Status      string         `json:"status"`
	Arguments   []byte         `json:"arguments"`
	Result      []byte         `json:"result"`
	LastError   sql.NullString `json:"last_error"`
	RetryCount  int            `json:"retry_count"`
	Options     jobs.Options   `json:"options"`
	ScheduledAt sql.NullTime   `json:"scheduled_at"`
	StartedAt   sql.NullTime   `json:"started_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (j job) ToJob() jobs.Job {
	return jobs.Job{
		ID:          j.ID,
		Queue:       j.Queue,
		Status:      j.Status,
		Arguments:   j.Arguments,
		Result:      j.Result,
		LastError:   j.LastError.String,
		RetryCount:  j.RetryCount,
		Options:     j.Options,
		ScheduledAt: j.ScheduledAt.Time,
		StartedAt:   j.StartedAt.Time,
		CreatedAt:   j.CreatedAt,
		UpdatedAt:   j.UpdatedAt,
	}
}

func (j *job) ScanDestinations() []interface{} {
	return []interface{}{
		&j.ID,
		&j.Status,
		&j.Queue,
		&j.Arguments,
		&j.Result,
		&j.LastError,
		&j.RetryCount,
		&j.Options,
		&j.CreatedAt,
		&j.UpdatedAt,
		&j.ScheduledAt,
	}
}
