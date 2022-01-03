package main

import (
	"database/sql"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.mpi-internal.com/SCM-Italy/gork/client"
)

type args struct {
	Wow int `json:"wow"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}

	for {
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		if err := client.Enqueue(tx, "default", args{Wow: rand.Int() % 200}); err != nil {
			panic(err)
		}

		if err := tx.Commit(); err != nil {
			panic(err)
		}

		time.Sleep(500 * time.Millisecond)
	}

}
