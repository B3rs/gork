package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/B3rs/gork/jobs"
	"github.com/B3rs/gork/workers"
	_ "github.com/lib/pq"
)

func main() {

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}

	pool := workers.NewWorkerPool(db,
		workers.WithGracefulShutdown(),
		workers.WithAdminUI(db, ":8080"),
	)
	pool.RegisterWorker("increase", IncreaseWorker{}, workers.WithInstances(3), workers.WithTimeout(10*time.Second)) // worker can be a struct method (so you can inject dependencies)
	pool.RegisterWorkerFunc("decrease", Decrease, workers.WithInstances(2))                                          // or a simple function

	if err := pool.Start(); err != nil {
		panic(err)
	}
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
