package workers

import (
	"time"
)

type workerConfig struct {
	worker    Worker
	instances int
	timeout   time.Duration
}

func newRegister() register {
	return register{}
}

type register map[string]workerConfig

// RegisterWorker registers a worker with the given name.
func (r register) RegisterWorker(queueName string, worker Worker, instances int, opts ...WorkerOptionFunc) {
	w := workerConfig{
		worker:    worker,
		instances: instances,
	}

	options := append(defaultWorkerOptions, opts...)

	for _, opt := range options {
		w = opt(w)
	}
	r[queueName] = w
}

// RegisterWorkerFunc registers a worker function with the given queue name.
func (r register) RegisterWorkerFunc(queueName string, worker WorkerFunc, instances int, opts ...WorkerOptionFunc) {

	r.RegisterWorker(queueName, funcWorker{f: worker}, instances, opts...)
}

func (r register) getWorkers() map[string]workerConfig {
	return r
}
