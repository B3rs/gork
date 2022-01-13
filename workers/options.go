package workers

import "time"

var (
	defaultPoolOptions = []PoolOptionFunc{
		WithSchedulerInterval(1 * time.Second),
		WithReaperInterval(10 * time.Second),
		WithLogger(&nopLogger{}),
	}

	defaultWorkerOptions = []WorkerOptionFunc{
		WithTimeout(1 * time.Minute),
	}
)

type PoolOptionFunc func(p *WorkerPool) *WorkerPool

func WithSchedulerInterval(interval time.Duration) PoolOptionFunc {
	return func(p *WorkerPool) *WorkerPool {
		p.schedulerSleepInterval = interval
		return p
	}
}

func WithReaperInterval(interval time.Duration) PoolOptionFunc {
	return func(p *WorkerPool) *WorkerPool {
		p.reaperInterval = interval
		return p
	}
}

func WithLogger(l logger) PoolOptionFunc {
	return func(p *WorkerPool) *WorkerPool {
		p.logger = l
		return p
	}
}

type WorkerOptionFunc func(workerConfig) workerConfig

func WithTimeout(d time.Duration) WorkerOptionFunc {
	return func(w workerConfig) workerConfig {
		w.timeout = d
		return w
	}
}
