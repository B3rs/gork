package jobs

import (
	"database/sql"
	"strings"
	"time"
)

type jobRow struct {
	ID          string         `json:"id"`
	Queue       string         `json:"queue"`
	Status      string         `json:"status"`
	Arguments   []byte         `json:"arguments"`
	Result      []byte         `json:"result"`
	LastError   sql.NullString `json:"last_error"`
	RetryCount  int            `json:"retry_count"`
	Options     Options        `json:"options"`
	ScheduledAt sql.NullTime   `json:"scheduled_at"`
	StartedAt   sql.NullTime   `json:"started_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (j jobRow) toJob() *Job {
	return &Job{
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

func (j jobRow) columns() []string {
	return []string{
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
}

func (j jobRow) stringColumns() string {
	return strings.Join(j.columns(), ", ")
}

func (j *jobRow) scanDestinations() []interface{} {
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
