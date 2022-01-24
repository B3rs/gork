package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/B3rs/gork/jobs"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTxWrapper_WrapTx(t *testing.T) {
	type args struct {
		ctx context.Context
		f   func(ctx context.Context, tx *sql.Tx) (interface{}, error)
	}
	tests := []struct {
		name            string
		args            args
		sqlExpectations func(mock sqlmock.Sqlmock)
		want            interface{}
		wantErr         error
	}{
		{
			name: "happy path",
			args: args{
				ctx: context.Background(),
				f: func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
					_, _ = tx.Exec("INSERT INTO jobs")
					return "result", nil
				},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO jobs").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			want:    "result",
			wantErr: nil,
		},
		{
			name: "begin error should be returned",
			args: args{
				ctx: context.Background(),
				f: func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
					_, _ = tx.Exec("INSERT INTO jobs")
					return "result", nil
				},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			wantErr: errors.New("begin error"),
		},
		{
			name: "function error should trigger rollback",
			args: args{
				ctx: context.Background(),
				f: func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
					return nil, errors.New("function error")
				},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			wantErr: errors.New("function error"),
		},
		{
			name: "job not found should be committed",
			args: args{
				ctx: context.Background(),
				f: func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
					return nil, jobs.ErrJobNotFound
				},
			},
			sqlExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			want:    nil,
			wantErr: jobs.ErrJobNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err, "an error '%s' was not expected when opening a stub database connection", err)
			defer db.Close()

			tt.sqlExpectations(mock)

			w := TxWrapper{
				db: db,
			}
			got, err := w.WrapTx(tt.args.ctx, tt.args.f)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)

			// we make sure that all expectations were met
			err = mock.ExpectationsWereMet()
			assert.Nil(t, err, "there were unfulfilled expectations: %s", err)
		})
	}
}
