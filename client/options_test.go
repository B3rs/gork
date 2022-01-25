package client

import (
	"testing"
	"time"

	"github.com/B3rs/gork/jobs"
	"github.com/stretchr/testify/assert"
)

func TestWithMaxRetries(t *testing.T) {
	job := jobs.Job{}
	job = WithMaxRetries(3)(job)
	assert.Equal(t, 3, job.Options.MaxRetries)
}

func TestWithRetryInterval(t *testing.T) {
	job := jobs.Job{}
	job = WithRetryInterval(time.Second)(job)
	assert.Equal(t, time.Second, job.Options.RetryInterval)
}

func TestWithScheduleTime(t *testing.T) {
	job := jobs.Job{}
	job = WithRetryInterval(time.Second)(job)
	assert.Equal(t, time.Second, job.Options.RetryInterval)
}

func Test_scheduleImmediately(t *testing.T) {
	job := jobs.Job{}
	now = func() time.Time { return time.Time{}.Add(time.Hour) }
	job = scheduleImmediately()(job)
	assert.Equal(t, time.Time{}.Add(time.Hour), job.ScheduledAt)
}
