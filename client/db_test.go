package client

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/B3rs/gork/jobs"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type noopWrapper struct {
}

func (n noopWrapper) WrapTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) (interface{}, error)) (interface{}, error) {
	return f(ctx, nil)
}

func TestClient_Schedule(t *testing.T) {
	later := time.Now().Add(time.Hour)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txmock := NewMockTxClient(ctrl)

	txmock.EXPECT().Schedule(context.Background(), "jobid", "awesomequeue", "7", gomock.Any()).Return(nil).Times(1)
	c := Client{
		txClientFactory: func(tx *sql.Tx) TxClient {
			return txmock
		},
		txWrapper: noopWrapper{},
	}

	err := c.Schedule(context.Background(), "jobid", "awesomequeue", "7", WithScheduleTime(later))
	assert.NoError(t, err)
}

func TestClient_Schedule_error(t *testing.T) {
	later := time.Now().Add(time.Hour)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockTxClient(ctrl)

	txmock.EXPECT().Schedule(context.Background(), "jobid", "awesomequeue", "7", gomock.Any()).Return(errors.New("error")).Times(1)

	c := Client{
		txClientFactory: func(tx *sql.Tx) TxClient {
			return txmock
		},
		txWrapper: noopWrapper{},
	}

	err := c.Schedule(context.Background(), "jobid", "awesomequeue", "7", WithScheduleTime(later))
	assert.Error(t, err, "should return error if error occurs")
}

func TestClient_Cancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockTxClient(ctrl)

	txmock.EXPECT().Cancel(gomock.Any(), "jobid").Return(nil).Times(1)

	c := Client{
		txClientFactory: func(tx *sql.Tx) TxClient {
			return txmock
		},
		txWrapper: noopWrapper{},
	}

	err := c.Cancel(context.Background(), "jobid")
	assert.NoError(t, err)
}

func TestClient_ScheduleNow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockTxClient(ctrl)

	txmock.EXPECT().ScheduleNow(gomock.Any(), "jobid").Return(nil).Times(1)

	c := Client{
		txClientFactory: func(tx *sql.Tx) TxClient {
			return txmock
		},
		txWrapper: noopWrapper{},
	}

	err := c.ScheduleNow(context.Background(), "jobid")
	assert.NoError(t, err)
}

func TestClient_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockTxClient(ctrl)

	txmock.EXPECT().Get(gomock.Any(), "jobid").Return(jobs.Job{
		ID:        "jobid",
		Queue:     "awesomequeue",
		Arguments: json.RawMessage(`"7"`),
		Status:    jobs.StatusScheduled,
		Options: jobs.Options{
			MaxRetries:    0,
			RetryInterval: 10 * time.Second,
		},
	}, nil).Times(1)

	c := Client{
		txClientFactory: func(tx *sql.Tx) TxClient {
			return txmock
		},
		txWrapper: noopWrapper{},
	}

	got, err := c.Get(context.Background(), "jobid")
	assert.NoError(t, err)
	assert.Equal(t, jobs.Job{
		ID:        "jobid",
		Queue:     "awesomequeue",
		Arguments: json.RawMessage(`"7"`),
		Status:    jobs.StatusScheduled,
		Options: jobs.Options{
			MaxRetries:    0,
			RetryInterval: 10 * time.Second,
		},
	}, got)

}

func TestClient_Get_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockTxClient(ctrl)

	txmock.EXPECT().Get(gomock.Any(), "jobid").Return(jobs.Job{}, errors.New("error")).Times(1)

	c := Client{
		txClientFactory: func(tx *sql.Tx) TxClient {
			return txmock
		},
		txWrapper: noopWrapper{},
	}

	got, err := c.Get(context.Background(), "jobid")
	assert.Error(t, err)
	assert.Equal(t, jobs.Job{}, got)

}

func TestClient_WithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockTxClient(ctrl)

	c := Client{
		txClientFactory: func(tx *sql.Tx) TxClient {
			return txmock
		},
		txWrapper: noopWrapper{},
	}

	got := c.WithTx(nil)
	assert.Equal(t, txmock, got)

}
