package workers

import (
	"context"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

type workerInfo struct {
	worker    Worker
	instances int
}

func newRegister() register {
	return register{}
}

type register map[string]workerInfo

// RegisterWorker registers a worker with the given name.
func (r register) RegisterWorker(queueName string, worker Worker, instances int) {
	r[queueName] = workerInfo{
		worker:    worker,
		instances: instances,
	}
}

// funcWorker is just a wrapper for a WorkerFunc.
type funcWorker struct {
	f WorkerFunc
}

// Execute runs the worker function for the given job.
func (f funcWorker) Execute(ctx context.Context, job jobs.Job) (interface{}, error) {
	return f.f(ctx, job)
}

// RegisterWorkerFunc registers a worker function with the given queue name.
func (r register) RegisterWorkerFunc(queueName string, worker WorkerFunc, instances int) {

	r[queueName] = workerInfo{
		worker:    funcWorker{f: worker},
		instances: instances,
	}
}

func (r register) getWorkers() map[string]workerInfo {
	return r
}
