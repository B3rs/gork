package main

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/B3rs/gork/client"
	_ "github.com/lib/pq"
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

	c := client.NewDBClient(db)

	for i := 0; i < 10; i++ {

		if err := c.Schedule(context.Background(), "increase_"+strconv.Itoa(i), "increase", args{Wow: rand.Int() % 200}); err != nil {
			panic(err)
		}

		if err := c.Schedule(context.Background(), "decrease_"+strconv.Itoa(i), "decrease", args2{Bau: rand.Int() % 200}); err != nil {
			panic(err)
		}
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	txc := client.NewTxClient(tx)

	if err := txc.Cancel(context.Background(), "increase_3"); err != nil {
		panic(err)
	}

	if err := txc.Cancel(context.Background(), "increase_6"); err != nil {
		panic(err)
	}

	if err := txc.Cancel(context.Background(), "increase_9"); err != nil {
		panic(err)
	}

	if err := txc.Schedule(context.Background(), "1218", "increase", args{Wow: 123}, client.WithMaxRetries(3), client.WithRetryInterval(1*time.Second), client.WithScheduleTime(time.Now().Add(2*time.Second))); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

}
