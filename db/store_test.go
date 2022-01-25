package db

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/B3rs/gork/jobs"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStore_Update(t *testing.T) {
	now := time.Now()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE jobs
	SET
		status=\$1, 
		result=\$2, 
		last_error=\$3, 
		retry_count=\$4,
		scheduled_at=\$5,
		updated_at=now\(\)
	WHERE id = \$6`).WithArgs(
		jobs.StatusScheduled,
		json.RawMessage(`{}`),
		"",
		1,
		now,
		"1",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewStore(db)
	err = store.Update(context.TODO(), jobs.Job{
		ID:          "1",
		Status:      jobs.StatusScheduled,
		Result:      json.RawMessage(`{}`),
		LastError:   "",
		RetryCount:  1,
		ScheduledAt: now,
	})

	assert.Nil(t, err, "update should not return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)

}

func TestStore_Update_error(t *testing.T) {
	now := time.Now()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE jobs
	SET
		status=\$1, 
		result=\$2, 
		last_error=\$3, 
		retry_count=\$4,
		scheduled_at=\$5,
		updated_at=now\(\)
	WHERE id = \$6`).WithArgs(
		jobs.StatusScheduled,
		json.RawMessage(`{}`),
		"",
		1,
		now,
		"1",
	).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	store := NewStore(db)
	err = store.Update(context.TODO(), jobs.Job{
		ID:          "1",
		Status:      jobs.StatusScheduled,
		Result:      json.RawMessage(`{}`),
		LastError:   "",
		RetryCount:  1,
		ScheduledAt: now,
	})

	assert.Equal(t, errors.New("error"), err, "update should return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)

}

func TestStore_Create(t *testing.T) {
	now := time.Now()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO jobs \(id, queue, status, arguments, options, scheduled_at\) 
	VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).WithArgs(
		"1",
		"q",
		jobs.StatusScheduled,
		json.RawMessage(`{}`),
		jobs.Options{
			MaxRetries: 1,
		},
		now,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewStore(db)
	err = store.Create(context.TODO(), jobs.Job{
		ID:        "1",
		Queue:     "q",
		Status:    jobs.StatusScheduled,
		Arguments: json.RawMessage(`{}`),
		Options: jobs.Options{
			MaxRetries: 1,
		},
		ScheduledAt: now,
	})

	assert.Nil(t, err, "create should not return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_Create_error(t *testing.T) {
	now := time.Now()

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO jobs \(id, queue, status, arguments, options, scheduled_at\) 
	VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).WithArgs(
		"1",
		"q",
		jobs.StatusScheduled,
		json.RawMessage(`{}`),
		jobs.Options{
			MaxRetries: 1,
		},
		now,
	).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	store := NewStore(db)
	err = store.Create(context.TODO(), jobs.Job{
		ID:        "1",
		Queue:     "q",
		Status:    jobs.StatusScheduled,
		Arguments: json.RawMessage(`{}`),
		Options: jobs.Options{
			MaxRetries: 1,
		},
		ScheduledAt: now,
	})

	assert.Equal(t, errors.New("error"), err, "create should return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_Deschedule(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE jobs 
	SET 
		updated_at=now\(\), 
		status=\$1 
	WHERE 
		id = \$2 AND 
		status = \$3`).WithArgs(
		jobs.StatusCanceled,
		"1",
		jobs.StatusScheduled,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewStore(db)
	err = store.Deschedule(context.TODO(), "1")
	assert.Nil(t, err, "deschedule should not return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_Deschedule_error(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE jobs 
	SET 
		updated_at=now\(\), 
		status=\$1 
	WHERE 
		id = \$2 AND 
		status = \$3`).WithArgs(
		jobs.StatusCanceled,
		"1",
		jobs.StatusScheduled,
	).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	store := NewStore(db)
	err = store.Deschedule(context.TODO(), "1")
	assert.Equal(t, errors.New("error"), err, "deschedule should return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_ScheduleNow(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE jobs 
	SET 
		updated_at=now\(\), 
		scheduled_at=now\(\), 
		status=\$1 
	WHERE 
		id = \$2`).WithArgs(
		jobs.StatusScheduled,
		"1",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewStore(db)
	err = store.ScheduleNow(context.TODO(), "1")
	assert.Nil(t, err, "schedule now should not return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_ScheduleNow_error(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE jobs 
	SET 
		updated_at=now\(\), 
		scheduled_at=now\(\), 
		status=\$1 
	WHERE 
		id = \$2`).WithArgs(
		jobs.StatusScheduled,
		"1",
	).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	store := NewStore(db)
	err = store.ScheduleNow(context.TODO(), "1")
	assert.Equal(t, errors.New("error"), err, "schedule now should return an error")

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_Search(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT `+jobStringColumns+`
	FROM jobs
	WHERE id LIKE '\%' \|\| \$1 \|\| '\%' 
	ORDER BY scheduled_at DESC 
	LIMIT \$2 OFFSET \$3`).WithArgs(
		"search",
		10,
		1,
	).WillReturnRows(sqlmock.NewRows(jobColumns).
		AddRow(jobQueryResult...).
		AddRow(jobQueryResult...),
	)
	mock.ExpectCommit()

	store := NewStore(db)
	got, err := store.Search(context.TODO(), 10, 1, "search")

	assert.Nil(t, err, "search should not return an error")
	assert.Equal(t, []jobs.Job{expectedJob, expectedJob}, got)

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_Search_empty(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT `+jobStringColumns+`
	FROM jobs
	ORDER BY scheduled_at DESC 
	LIMIT \$1 OFFSET \$2`).WithArgs(
		10,
		1,
	).WillReturnRows(sqlmock.NewRows(jobColumns).
		AddRow(jobQueryResult...).
		AddRow(jobQueryResult...),
	)
	mock.ExpectCommit()

	store := NewStore(db)
	got, err := store.Search(context.TODO(), 10, 1, "")

	assert.Nil(t, err, "search should not return an error")
	assert.Equal(t, []jobs.Job{expectedJob, expectedJob}, got)

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_Get(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT ` + jobStringColumns + `
	FROM jobs 
	WHERE id = \$1`).WithArgs(
		"1",
	).WillReturnRows(sqlmock.NewRows(jobColumns).
		AddRow(jobQueryResult...),
	)
	mock.ExpectCommit()

	store := NewStore(db)
	got, err := store.Get(context.TODO(), "1")

	assert.Nil(t, err, "search should not return an error")
	assert.Equal(t, expectedJob, got)

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}

func TestStore_Get_error(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT ` + jobStringColumns + `
	FROM jobs 
	WHERE id = \$1`).WithArgs(
		"1",
	).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	store := NewStore(db)
	got, err := store.Get(context.TODO(), "1")

	assert.Equal(t, errors.New("error"), err, "search should return an error")
	assert.Equal(t, jobs.Job{}, got)

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
}
