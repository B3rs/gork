package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.mpi-internal.com/SCM-Italy/gork/jobs"
	"github.mpi-internal.com/SCM-Italy/gork/workers"
)

func main() {

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}

	pool := workers.NewWorkerPool(db, 100*time.Millisecond)
	pool.RegisterWorker("increase", IncreaseWorker{}, 3) // worker can be a struct method (so you can inject dependencies)
	pool.RegisterWorkerFunc("decrease", Decrease, 2)     // or a simple function

	sigc := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		fmt.Println("\n\n\nReceived an interrupt, stopping services...\n\n")
		pool.Stop()
		done <- true
	}()

	pool.Start()
	<-done
}

type args struct {
	Wow int `json:"wow"`
}

type IncreaseWorker struct {
}

func (w IncreaseWorker) Execute(ctx context.Context, job jobs.Job) (interface{}, error) {

	fmt.Println("start increase job", job.ID, string(job.Arguments))
	a := args{}

	if err := job.ParseArguments(&a); err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(rand.Int()%2000) * time.Millisecond)

	if a.Wow == 123 {
		return nil, errors.New("error, number is 123")
	}

	result := a.Wow + 1
	return result, nil
}

type args2 struct {
	Bau int `json:"bau"`
}

func Decrease(ctx context.Context, job jobs.Job) (interface{}, error) {

	fmt.Println("start decrease job", job.ID, string(job.Arguments))
	a := args2{}

	if err := job.ParseArguments(&a); err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(rand.Int()%2000) * time.Millisecond)

	if a.Bau == 21 {
		return nil, errors.New("error, number is 21")
	}
	result := a.Bau - 1
	return result, nil
}
