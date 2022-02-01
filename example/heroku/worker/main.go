package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/B3rs/gork/jobs"
	"github.com/B3rs/gork/workers"
	_ "github.com/lib/pq"
)

func main() {

	// open a db connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	manageErr(err)

	// create a worker pool
	pool := workers.NewWorkerPool(
		db,
		workers.WithGracefulShutdown(), // Add graceful shutdown if you like to complete all jobs before exiting
		workers.WithAdminUI(":"+os.Getenv("PORT")), // Add admin UI if you like
	)
	// register a worker on increase queue
	pool.RegisterWorker(
		"increase",
		IncreaseWorker{
			logFunction: log.Println,
		}, // by providing a worker struct (so you can inject dependencies)
		workers.WithInstances(3),            // and spawning 3 worker routines
		workers.WithTimeout(10*time.Second), // with a 10 second timeout for each job
	)

	// you can also add workers with a function
	pool.RegisterWorkerFunc(
		"lowerize",
		Lowerize,                 // by prviding a worker function
		workers.WithInstances(2), // and spawning 2 worker routines
	) // or a simple function

	if err := pool.Start(); err != nil {
		panic(err)
	}
}

type IncreaseArgs struct {
	IncreaseThis int `json:"increase_this"`
}

type IncreaseResult struct {
	Increased int `json:"increased"`
}

type IncreaseWorker struct {
	workers.DefaultWorker

	logFunction func(...interface{})
}

func (w IncreaseWorker) Execute(ctx context.Context, job jobs.Job) (interface{}, error) {
	t := time.Now()
	w.logFunction("start job", job.ID, string(job.Arguments))
	defer func() { w.logFunction("end job", job.ID, "in", time.Since(t)) }()

	// Parse the job arguments into a struct
	args := IncreaseArgs{}
	if err := job.ParseArguments(&args); err != nil {
		return nil, err
	}

	// Do some work
	time.Sleep(time.Duration(rand.Int()%2000) * time.Millisecond)

	// Introduce bugs :)
	if args.IncreaseThis == 123 {
		return nil, errors.New("error, number is 123")
	}

	// Return the result to be saved in the database
	return IncreaseResult{Increased: args.IncreaseThis + 1}, nil
}

// We can add a custom callback for failure, so we can do something like:
func (w IncreaseWorker) OnFailure(ctx context.Context, job jobs.Job) error {
	log.Println("job failed", job.ID, job.LastError)
	return nil
}

type LowerizeArgs struct {
	LowerizeThis string `json:"lowerize_this"`
}
type LowerizeResult struct {
	Lowerized string `json:"lowerized"`
}

func Lowerize(ctx context.Context, job jobs.Job) (interface{}, error) {
	t := time.Now()
	log.Println("start job", job.ID, string(job.Arguments))
	defer func() { log.Println("end job", job.ID, "in", time.Since(t)) }()

	// Parse the job arguments into a struct
	args := LowerizeArgs{}
	if err := job.ParseArguments(&args); err != nil {
		return nil, err
	}

	// Do some work
	time.Sleep(time.Duration(rand.Int()%2000) * time.Millisecond)

	// Save the result
	return LowerizeResult{Lowerized: strings.ToLower(args.LowerizeThis)}, nil
}
func manageErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
