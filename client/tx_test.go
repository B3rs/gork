package client

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/B3rs/gork/jobs"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_txClient_Schedule(t *testing.T) {
	later := time.Now().Add(time.Hour)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockdbTx(ctrl)

	txmock.EXPECT().Create(gomock.Any(), jobs.Job{
		ID:          "jobid",
		Queue:       "awesomequeue",
		Arguments:   json.RawMessage(`"7"`),
		Status:      jobs.StatusScheduled,
		ScheduledAt: later,
		Options: jobs.Options{
			MaxRetries:    0,
			RetryInterval: 10 * time.Second,
		},
	}).Return(nil).Times(1)

	c := txClient{tx: txmock}

	err := c.Schedule(context.Background(), "jobid", "awesomequeue", "7", WithScheduleTime(later))
	assert.NoError(t, err)
}

func Test_txClient_Schedule_params_error(t *testing.T) {
	later := time.Now().Add(time.Hour)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockdbTx(ctrl)

	c := txClient{tx: txmock}

	err := c.Schedule(context.Background(), "jobid", "awesomequeue", make(chan struct{}), WithScheduleTime(later))
	assert.Error(t, err, "should return error if arguments are not serializable")
}

func Test_txClient_Schedule_error(t *testing.T) {
	later := time.Now().Add(time.Hour)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockdbTx(ctrl)

	txmock.EXPECT().Create(gomock.Any(), jobs.Job{
		ID:          "jobid",
		Queue:       "awesomequeue",
		Arguments:   json.RawMessage(`"7"`),
		Status:      jobs.StatusScheduled,
		ScheduledAt: later,
		Options: jobs.Options{
			MaxRetries:    0,
			RetryInterval: 10 * time.Second,
		},
	}).Return(errors.New("error")).Times(1)

	c := txClient{tx: txmock}

	err := c.Schedule(context.Background(), "jobid", "awesomequeue", "7", WithScheduleTime(later))
	assert.Error(t, err, "should return error if error occurs")
}

func Test_txClient_Cancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockdbTx(ctrl)

	txmock.EXPECT().Deschedule(gomock.Any(), "jobid").Return(nil).Times(1)

	c := txClient{tx: txmock}

	err := c.Cancel(context.Background(), "jobid")
	assert.NoError(t, err)
}

func Test_txClient_ScheduleNow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockdbTx(ctrl)

	txmock.EXPECT().ScheduleNow(gomock.Any(), "jobid").Return(nil).Times(1)

	c := txClient{tx: txmock}

	err := c.ScheduleNow(context.Background(), "jobid")
	assert.NoError(t, err)
}

func Test_txClient_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmock := NewMockdbTx(ctrl)

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

	c := txClient{tx: txmock}

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
