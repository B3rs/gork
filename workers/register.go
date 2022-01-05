package workers

import (
	"context"
	"errors"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

type workerInfo struct {
	worker    Worker
	instances int
}

func newRegister() register {
	return register{
		workers: make(map[string]workerInfo),
	}
}

type register struct {
	workers map[string]workerInfo
}

func (r *register) RegisterWorker(name string, worker Worker, instances int) {
	r.workers[name] = workerInfo{
		worker:    worker,
		instances: instances,
	}
}

type funcWorker struct {
	f WorkerFunc
}

func (f funcWorker) Execute(ctx context.Context, job jobs.Job) (interface{}, error) {
	return f.f(ctx, job)
}

func (r *register) RegisterWorkerFunc(name string, worker WorkerFunc, instances int) {

	r.workers[name] = workerInfo{
		worker:    funcWorker{f: worker},
		instances: instances,
	}
}

func (r *register) GetWorker(name string) (workerInfo, error) {
	if worker, ok := r.workers[name]; ok {
		return worker, nil
	}
	return workerInfo{}, errors.New("worker not found")
}

func (r *register) GetWorkers() map[string]workerInfo {
	return r.workers
}
