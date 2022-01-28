package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/B3rs/gork/jobs"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestQueue_Pop(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE jobs
	SET 
		status=\$1, 
		updated_at=now\(\),
		started_at=now\(\)
	WHERE 
		id = \(
			SELECT id
			FROM jobs 
			WHERE status = \$2
				AND scheduled_at <= now\(\)
				AND queue = \$3
			ORDER BY scheduled_at ASC 
			FOR UPDATE SKIP LOCKED
			LIMIT 1 
		\)
	RETURNING `+jobStringColumns).WithArgs(
		jobs.StatusInitialized,
		jobs.StatusScheduled,
		"test",
	).
		WillReturnRows(sqlmock.NewRows(jobColumns).
			AddRow(jobQueryResult...),
		)
	mock.ExpectCommit()

	q := NewQueue(db, "test")
	got, err := q.Pop(context.TODO())

	assert.Nil(t, err, "Pop should not return an error")
	assert.Equal(t, expectedJob, got)

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestQueue_Pop_error(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE jobs
	SET 
		status=\$1, 
		updated_at=now\(\),
		started_at=now\(\)
	WHERE 
		id = \(
			SELECT id
			FROM jobs 
			WHERE status = \$2
				AND scheduled_at <= now\(\)
				AND queue = \$3
			ORDER BY scheduled_at ASC 
			FOR UPDATE SKIP LOCKED
			LIMIT 1 
		\)
	RETURNING `+jobStringColumns).WithArgs(
		jobs.StatusInitialized,
		jobs.StatusScheduled,
		"test",
	).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	q := NewQueue(db, "test")
	got, err := q.Pop(context.TODO())

	assert.Equal(t, errors.New("error"), err, "Pop should  return an error")
	assert.Equal(t, jobs.Job{}, got)

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestQueue_RequeueTimedOutJobs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE jobs
	SET 
		status=\$1, 
		started_at=null,
		retry_count=retry_count\+1,
		updated_at=now\(\)
	WHERE 
		started_at < \$2 AND
		status = \$3 AND
		queue = \$4`,
	).
		WithArgs(
			jobs.StatusScheduled,
			time.Time{}.Add(-time.Second),
			jobs.StatusInitialized,
			"test",
		).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	q := NewQueue(db, "test")
	q.now = func() time.Time { return time.Time{} }

	err = q.RequeueTimedOutJobs(context.TODO(), time.Second)

	assert.Nil(t, err, "RequeueTimedOutJobs should not return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestQueue_RequeueTimedOutJobs_error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE jobs
	SET 
		status=\$1, 
		started_at=null,
		retry_count=retry_count\+1,
		updated_at=now\(\)
	WHERE 
		started_at < \$2 AND
		status = \$3 AND
		queue = \$4`,
	).
		WithArgs(
			jobs.StatusScheduled,
			time.Time{}.Add(-time.Second),
			jobs.StatusInitialized,
			"test",
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	q := NewQueue(db, "test")
	q.now = func() time.Time { return time.Time{} }

	err = q.RequeueTimedOutJobs(context.TODO(), time.Second)

	assert.Equal(t, errors.New("error"), err, "RequeueTimedOutJobs should return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}
