package jobs

import (
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
	Options     Options
	ScheduledAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (j Job) ParseArguments(dest interface{}) error {
	return json.Unmarshal(j.Arguments, dest)
}

// SetStatus sets the status of the job
func (j *Job) SetStatus(status string) {
	j.Status = status
}

// SetLastError sets the last error of the job
func (j *Job) SetLastError(err error) {
	j.LastError = err.Error()
}

// SetResult sets the result of the job
func (j *Job) SetResult(res interface{}) error {
	if res == nil {
		return nil
	}
	encoded, err := json.Marshal(res)
	if err != nil {
		return err
	}
	j.Result = encoded
	return nil
}
