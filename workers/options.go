package workers

import "time"

var (
	defaultPoolOptions = []PoolOptionFunc{
		WithSchedulerInterval(1 * time.Second),
		WithReaperInterval(10 * time.Second),
	}

	defaultWorkerOptions = []WorkerOptionFunc{
		WithTimeout(1 * time.Minute),
	}
)

type PoolOptionFunc func(p poolConfig) poolConfig

func WithSchedulerInterval(interval time.Duration) PoolOptionFunc {
	return func(c poolConfig) poolConfig {
		c.schedulerSleepInterval = interval
		return c
	}
}

func WithReaperInterval(interval time.Duration) PoolOptionFunc {
	return func(c poolConfig) poolConfig {
		c.reaperInterval = interval
		return c
	}
}

type WorkerOptionFunc func(workerConfig) workerConfig

func WithTimeout(d time.Duration) WorkerOptionFunc {
	return func(w workerConfig) workerConfig {
		w.timeout = d
		return w
	}
}
