package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/B3rs/gork/jobs"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	now            = time.Now()
	jobQueryResult = []driver.Value{"jobid", "scheduled", "jobq", []byte(`{"arg1":"val1"}`), []byte(`{"res1":"resval1"}`), nil, 0, []byte(`{}`), now, now, now}
	expectedJob    = jobs.Job{
		ID:          "jobid",
		Status:      "scheduled",
		Queue:       "jobq",
		Arguments:   json.RawMessage(`{"arg1":"val1"}`),
		Result:      json.RawMessage(`{"res1":"resval1"}`),
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
)

func Test_queryJob(t *testing.T) {

	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []struct {
		name            string
		args            args
		sqlExpectations func(mock sqlmock.Sqlmock)
		want            jobs.Job
		wantErr         error
	}{
		{
			name: "happy path",
			args: args{
				ctx:   context.Background(),
				query: `SELECT * FROM jobs WHERE id = $1`,
				args:  []interface{}{"1"},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM jobs WHERE id = \$1`).WillReturnRows(
					sqlmock.NewRows(jobColumns).AddRow(jobQueryResult...),
				)
			},
			want: expectedJob,
		},
		{
			name: "job not found",
			args: args{
				ctx:   context.Background(),
				query: `SELECT * FROM jobs WHERE id = $1`,
				args:  []interface{}{"1"},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM jobs WHERE id = \$1`).WillReturnError(sql.ErrNoRows)
			},
			wantErr: jobs.ErrJobNotFound,
		},
		{
			name: "query error",
			args: args{
				ctx:   context.Background(),
				query: `SELECT * FROM jobs WHERE id = $1`,
				args:  []interface{}{"1"},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM jobs WHERE id = \$1`).WillReturnError(errors.New("query error"))
			},
			wantErr: errors.New("query error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
			defer db.Close()

			mock.ExpectBegin()
			tt.sqlExpectations(mock)

			tx, err := db.Begin()
			assert.Nil(t, err)

			got, err := queryJob(tt.args.ctx, tx, tt.args.query, tt.args.args...)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)

			// we make sure that all expectations were met
			err = mock.ExpectationsWereMet()
			assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
		})
	}
}

func Test_queryJobs(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []struct {
		name            string
		args            args
		sqlExpectations func(mock sqlmock.Sqlmock)
		want            []jobs.Job
		wantErr         error
	}{
		{
			name: "happy path",
			args: args{
				ctx:   context.Background(),
				query: `SELECT * FROM jobs`,
				args:  []interface{}{"1"},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM jobs`).WillReturnRows(
					sqlmock.NewRows(jobColumns).
						AddRow(jobQueryResult...).
						AddRow(jobQueryResult...),
				)
			},
			want: []jobs.Job{expectedJob, expectedJob},
		},
		{
			name: "query error",
			args: args{
				ctx:   context.Background(),
				query: `SELECT * FROM jobs`,
				args:  []interface{}{"1"},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM jobs`).WillReturnError(errors.New("query error"))
			},
			wantErr: errors.New("query error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
			defer db.Close()

			mock.ExpectBegin()
			tt.sqlExpectations(mock)

			tx, err := db.Begin()
			assert.Nil(t, err)

			got, err := queryJobs(tt.args.ctx, tx, tt.args.query, tt.args.args...)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)

			// we make sure that all expectations were met
			err = mock.ExpectationsWereMet()
			assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
		})
	}
}

func Test_exec(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO jobs`).WillReturnError(errors.New("exec error"))

	tx, err := db.Begin()
	assert.Nil(t, err)

	err = exec(context.Background(), tx, `INSERT INTO jobs`, "arg1")

	assert.Equal(t, errors.New("exec error"), err)

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "there were unfulfilled expectations: %s", err)

}
