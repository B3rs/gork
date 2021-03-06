package client

import (
	"time"

	"github.com/B3rs/gork/jobs"
)

var (
	now = time.Now

	defaultOptions = []OptionFunc{
		WithRetryInterval(10 * time.Second),
		scheduleImmediately(),
	}
)

type OptionFunc func(j jobs.Job) jobs.Job

func WithMaxRetries(tries int) OptionFunc {
	return func(j jobs.Job) jobs.Job {
		j.Options.MaxRetries = tries
		return j
	}
}

func WithRetryInterval(interval time.Duration) OptionFunc {
	return func(j jobs.Job) jobs.Job {
		j.Options.RetryInterval = interval
		return j
	}
}

func WithScheduleTime(t time.Time) OptionFunc {
	return func(j jobs.Job) jobs.Job {
		j.ScheduledAt = t
		return j
	}
}

func scheduleImmediately() OptionFunc {
	return func(j jobs.Job) jobs.Job {
		j.ScheduledAt = now()
		return j
	}
}
