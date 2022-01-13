package workers

import (
	"context"
	"testing"

	"github.com/B3rs/gork/jobs"
	"github.com/stretchr/testify/assert"
)

type testWorker struct {
	res interface{}
	err error
}

func (w testWorker) Execute(context.Context, jobs.Job) (interface{}, error) {
	return w.res, w.err
}

func Test_register_RegisterWorker(t *testing.T) {

	r := newRegister()
	r.RegisterWorker("queue", testWorker{res: "best worker"}, 7)

	assert.Len(t, r.getWorkers(), 1)

	res, err := r["queue"].worker.Execute(context.Background(), jobs.Job{})
	assert.Equal(t, "best worker", res)
	assert.Nil(t, err)

}

func Test_register_RegisterWorkerFunc(t *testing.T) {

	r := newRegister()
	r.RegisterWorkerFunc("queue", func(ctx context.Context, job jobs.Job) (interface{}, error) {
		return "best worker", nil
	}, 7)

	assert.Len(t, r.getWorkers(), 1)

	res, err := r["queue"].worker.Execute(context.Background(), jobs.Job{})
	assert.Equal(t, "best worker", res)
	assert.Nil(t, err)
}
