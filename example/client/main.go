package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/B3rs/gork/client"
	_ "github.com/lib/pq"
)

type IncreaseArgs struct {
	IncreaseThis int `json:"increase_this"`
}

type LowerizeArgs struct {
	LowerizeThis string `json:"lowerize_this"`
}

func main() {
	// normally open a db connection with the standard sql package
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}

	// create a client
	c := client.NewClient(db)

	// schedule a job
	err = c.Schedule(
		context.Background(),
		"increase_1",
		"increase",
		IncreaseArgs{IncreaseThis: 1234},
	)
	manageErr(err)

	// this job is likely to fail so we can add some retries
	err = c.Schedule(
		context.Background(),
		"increase_2",
		"increase",
		IncreaseArgs{IncreaseThis: 123},
		client.WithMaxRetries(3),                // retry 3 times
		client.WithRetryInterval(2*time.Second), // wait 2 seconds between retries
	)
	manageErr(err)

	// we can also schedule a job that will be executed in the future
	desiredExecutionTime := time.Now().Add(5 * time.Second)
	err = c.Schedule(
		context.Background(),
		"increase_3",
		"increase",
		IncreaseArgs{IncreaseThis: 456},
		client.WithScheduleTime(desiredExecutionTime), // schedule in 5 seconds
	)
	manageErr(err)

	// we can think about scheduling a job in the future
	desiredExecutionTime = time.Now().Add(5 * time.Second)
	err = c.Schedule(
		context.Background(),
		"increase_4",
		"increase",
		IncreaseArgs{IncreaseThis: 658},
		client.WithScheduleTime(desiredExecutionTime), // schedule in 5 seconds
	)
	manageErr(err)

	// and cancel it if we want
	err = c.Cancel(context.Background(), "increase_4")
	manageErr(err)

	// we can schedule different kind of jobs
	err = c.Schedule(
		context.Background(),
		"lowerize_woof",
		"lowerize",                         // by using a different queue
		LowerizeArgs{LowerizeThis: "WoOf"}, // With different params
	)
	manageErr(err)

	// we can start a standard sql transaction
	tx, err := db.Begin()
	manageErr(err)

	// and do everything we did before inside a transaction
	err = c.WithTx(tx).Schedule(
		context.Background(),
		"lowerize_meow",
		"lowerize",
		LowerizeArgs{LowerizeThis: "MeOOOw"},
		client.WithMaxRetries(3),
		client.WithRetryInterval(1*time.Second),
		client.WithScheduleTime(time.Now().Add(2*time.Second)),
	)
	manageErr(err)

	// and commit the transaction atomically when all database operations are done
	err = tx.Commit()
	manageErr(err)
}

func manageErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
