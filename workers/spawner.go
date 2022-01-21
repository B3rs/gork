package workers

import (
	"context"
	"sync"
)

func newSpawner(ctx context.Context, errChan chan<- error) spawner {
	ctx, cancel := context.WithCancel(ctx)
	return spawner{
		wg:       &sync.WaitGroup{},
		ctx:      ctx,
		shutdown: cancel,
		errChan:  errChan,
	}
}

// spawner keeps track of running routines
type spawner struct {
	wg *sync.WaitGroup

	ctx      context.Context
	shutdown func()
	errChan  chan<- error
}

//go:generate mockgen -destination=./runner_mock.go -package=workers -source=spawner.go
type runner interface {
	Run(context.Context, chan<- error)
}

func (s spawner) Spawn(r runner) {
	s.wg.Add(1)
	go func() {
		r.Run(s.ctx, s.errChan)
		s.wg.Done()
	}()
}

func (s spawner) Wait() {
	s.wg.Wait()
}

func (s spawner) Shutdown() {
	s.shutdown()
}
