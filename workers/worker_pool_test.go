package workers

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testErrHandler struct {
	errors []error
}

func (h *testErrHandler) Handle(err error) {
	h.errors = append(h.errors, err)
}

func TestWorkerPool_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := NewMockSpawner(ctrl)

	errChan := make(chan error, 3)

	errChan <- errors.New("error1")
	errChan <- errors.New("error2")
	errChan <- errors.New("error3")

	errHandler := &testErrHandler{}

	_, cancel := context.WithCancel(context.Background())
	w := &WorkerPool{
		register:       newRegister(),
		spawner:        s,
		errChan:        errChan,
		shutdown:       cancel,
		reaperInterval: 1,
		queueFactory: func(name string) Queue {
			return nil
		},
		errorHandler: errHandler.Handle,
	}

	w.RegisterWorkerFunc("test", nil, 3)
	w.RegisterWorkerFunc("test2", nil, 2)

	// 3 workers + 1 reaper and 2 workers + 1 reaper
	s.EXPECT().Spawn(gomock.Any()).Times(7)
	s.EXPECT().Wait().Times(1)

	w.Start()

	assert.Equal(t, []error{errors.New("error1"), errors.New("error2"), errors.New("error3")}, errHandler.errors, "should wait for errors to be processed")
}

func TestWorkerPool_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := NewMockSpawner(ctrl)
	s.EXPECT().Shutdown()

	w := WorkerPool{
		spawner: s,
	}
	w.Stop()
}
