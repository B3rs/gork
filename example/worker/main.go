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

type args struct {
	Wow int `json:"wow"`
}

func main() {

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}

	worker := workers.NewWorker(db, "default", workFunc)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		worker.Stop()
	}()

	panic(worker.Start(10))

}

// this stuff can improve dramatically with generics
func workFunc(ctx context.Context, job jobs.Job) (interface{}, error) {

	fmt.Println("start processing job", job.ID)
	a := args{}
	err := job.ParseArguments(&a)
	if err != nil {
		return nil, err
	}
	fmt.Printf("arguments %v\n", a)

	time.Sleep(time.Duration(rand.Int()%2000) * time.Millisecond)

	if a.Wow == 123 {
		return nil, errors.New("error, number is 123")
	}
	return a.Wow + 1, nil
}
