package client

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

func Enqueue(tx *sql.Tx, queueName string, arguments interface{}) error {

	encoded, err := json.Marshal(arguments)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO jobs (queue, status, arguments) VALUES ($1, $2, $3)", queueName, jobs.StatusScheduled, encoded)
	return err
}

func ScheduleAt(tx *sql.Tx, queueName string, arguments interface{}, scheduledAt time.Time) error {

	encoded, err := json.Marshal(arguments)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO jobs (queue, status, arguments, scheduled_at) VALUES ($1, $2, $3, $4)", queueName, jobs.StatusScheduled, encoded, scheduledAt)
	return err
}
