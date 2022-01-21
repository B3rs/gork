package workers

import "time"

var (
	defaultPoolOptions = []PoolOptionFunc{
		WithSchedulerInterval(1 * time.Second),
		WithReaperInterval(10 * time.Second),
		WithErrorHandler(defaultErrorHandler),
	}

	defaultWorkerOptions = []WorkerOptionFunc{
		WithTimeout(1 * time.Minute),
		WithInstances(1),
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

func WithErrorHandler(f func(error)) PoolOptionFunc {
	return func(p *WorkerPool) *WorkerPool {
		p.errorHandler = f
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

func WithInstances(i int) WorkerOptionFunc {
	return func(w workerConfig) workerConfig {
		w.instances = i
		return w
	}
}
