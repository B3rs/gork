package workers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/B3rs/gork/web"
)

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

func WithGracefulShutdown() PoolOptionFunc {
	return func(p *WorkerPool) *WorkerPool {
		routine := func() error {
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc,
				syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT)

			select {
			case <-sigc:
			case <-p.spawner.Done():
			}

			fmt.Printf("\n\n\nReceived an interrupt, stopping services...\n\n")

			p.Stop()
			return nil
		}

		p.coRoutines = append(p.coRoutines, routine)
		return p
	}
}

func WithAdminUI(addr string) PoolOptionFunc {
	return func(p *WorkerPool) *WorkerPool {
		routine := func() error {
			s := web.NewServer(p.store)
			// Start server
			go func() {

				quit := make(chan os.Signal, 1)
				signal.Notify(quit,
					syscall.SIGHUP,
					syscall.SIGINT,
					syscall.SIGTERM,
					syscall.SIGQUIT)

				select {
				case <-quit:
				case <-p.spawner.Done():
				}

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_ = s.Shutdown(ctx)

			}()
			return s.Start(addr)
		}

		p.coRoutines = append(p.coRoutines, routine)
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
