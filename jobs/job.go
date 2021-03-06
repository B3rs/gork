package jobs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	// StatusInitialized is the status of a job that is currently being processed
	StatusInitialized = "initialized"
)

var (
	// // ErrNoJobsAvailable is returned when no jobs are available
	// ErrNoJobsAvailable = errors.New("no jobs available")
	// ErrJobNotFound is returned when no job is found
	ErrJobNotFound = errors.New("job not found")
)

// Job could use generics for params and result
type Job struct {
	ID          string          `json:"id"`
	Queue       string          `json:"queue"`
	Status      string          `json:"status"`
	Arguments   json.RawMessage `json:"arguments"`
	Result      json.RawMessage `json:"result"`
	LastError   string          `json:"last_error"`
	RetryCount  int             `json:"retry_count"`
	Options     Options         `json:"options"`
	ScheduledAt time.Time       `json:"scheduled_at,omitempty"`
	StartedAt   time.Time       `json:"started_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (j Job) ParseArguments(dest interface{}) error {
	return json.Unmarshal([]byte(j.Arguments), dest)
}

func (j Job) ShouldRetry() bool {
	return j.RetryCount < j.Options.MaxRetries
}

func (j Job) ScheduleRetry(t time.Time) Job {
	j.RetryCount++
	j.ScheduledAt = t
	j.Status = StatusScheduled
	return j
}

// SetStatus sets the status of the job
func (j Job) SetStatus(status string) Job {
	j.Status = status
	return j
}

// SetLastError sets the last error of the job
func (j Job) SetLastError(err error) Job {
	j.LastError = err.Error()
	return j
}

// SetResult sets the result of the job
func (j Job) SetResult(res interface{}) (Job, error) {
	if res == nil {
		return j, nil
	}
	encoded, err := json.Marshal(res)
	if err != nil {
		return j, err
	}
	j.Result = encoded
	return j, nil
}

// SetArguments sets the arguments of the job
func (j Job) SetArguments(args interface{}) (Job, error) {
	if args == nil {
		return j, nil
	}
	encoded, err := json.Marshal(args)
	if err != nil {
		return j, err
	}
	j.Arguments = encoded
	return j, nil
}

type Options struct {
	MaxRetries    int           `json:"max_retries"`
	RetryInterval time.Duration `json:"retry_interval"`
}

func (o Options) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *Options) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), o)
}
