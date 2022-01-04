package client

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.mpi-internal.com/SCM-Italy/gork/jobs"
)

// Schedule schedules a job in the queue to be executed as soon as possible
func Schedule(tx *sql.Tx, id string, queueName string, arguments interface{}) error {

	encoded, err := json.Marshal(arguments)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO jobs (id, queue, status, arguments) VALUES ($1, $2, $3, $4)", id, queueName, jobs.StatusScheduled, encoded)
	return err
}

// ScheduleAt schedules a job in the queue to be executed at the given time
func ScheduleAt(tx *sql.Tx, id string, queueName string, arguments interface{}, scheduledAt time.Time) error {

	encoded, err := json.Marshal(arguments)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO jobs (id, queue, status, arguments, scheduled_at) VALUES ($1, $2, $3, $4, $5)", id, queueName, jobs.StatusScheduled, encoded, scheduledAt)
	return err
}

// ScheduleAfter schedules a job in the queue to be executed after the given duration
func ScheduleAfter(tx *sql.Tx, id string, queueName string, arguments interface{}, after time.Duration) error {
	return ScheduleAt(tx, id, queueName, arguments, time.Now().Add(after))
}

// Cancel cancels a job in the queue if not already executed
func Cancel(tx *sql.Tx, id string) error {
	_, err := tx.Exec("UPDATE jobs SET updated_at=now(), scheduled_at=NULL, status=$1 WHERE id = $2 AND status = $3", jobs.StatusCanceled, id, jobs.StatusScheduled)
	return err
}
