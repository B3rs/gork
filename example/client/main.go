package main

import (
	"database/sql"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.mpi-internal.com/SCM-Italy/gork/client"
)

type args struct {
	Wow int `json:"wow"`
}

type args2 struct {
	Bau int `json:"bau"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		if err := client.Schedule(tx, "increase_"+strconv.Itoa(i), "increase", args{Wow: rand.Int() % 200}); err != nil {
			panic(err)
		}

		if err := client.Schedule(tx, "decrease_"+strconv.Itoa(i), "decrease", args2{Bau: rand.Int() % 200}); err != nil {
			panic(err)
		}

		if err := tx.Commit(); err != nil {
			panic(err)
		}
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	if err := client.Cancel(tx, "3"); err != nil {
		panic(err)
	}

	if err := client.Cancel(tx, "6"); err != nil {
		panic(err)
	}

	if err := client.Cancel(tx, "9"); err != nil {
		panic(err)
	}

	if err := client.ScheduleAt(tx, "1218", "increase", args{Wow: 123}, time.Now().Add(10*time.Second)); err != nil {
		panic(err)
	}

	if err := client.ScheduleAfter(tx, "3654", "increase", args{Wow: 5684}, 5*time.Second); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

}
