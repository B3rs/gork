package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
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

	worker := workers.NewWorker(db, "default", func(ctx context.Context, job jobs.Job) (interface{}, error) {
		time.Sleep(time.Duration(rand.Int()%2000) * time.Millisecond)
		a := args{}
		err := job.ParseArguments(&a)
		if err != nil {
			return nil, err
		}
		fmt.Println("wow", a)

		if a.Wow == 123 {
			return nil, errors.New("error")
		}
		return args{Wow: a.Wow + 1}, nil
	})

	panic(worker.Start(10))

}
