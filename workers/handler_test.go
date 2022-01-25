package workers

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	time "time"

	"github.com/B3rs/gork/db"
	"github.com/B3rs/gork/jobs"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Handler_Handle(t *testing.T) {
	tests := []struct {
		name               string
		job                jobs.Job
		workerExpectation  func(w *MockWorker)
		updaterExpectation func(u *db.MockJobsStore)
		wantErr            error
	}{
		{
			name: "should set status completed and set a result if job execution succeeds",
			job:  jobs.Job{ID: "1"},
			workerExpectation: func(w *MockWorker) {
				w.EXPECT().Execute(gomock.Any(), gomock.Eq(jobs.Job{ID: "1"})).Return("resultstring", nil)
			},
			updaterExpectation: func(u *db.MockJobsStore) {
				b, _ := json.Marshal("resultstring")
				u.EXPECT().Update(gomock.Any(), gomock.Eq(jobs.Job{ID: "1", Status: jobs.StatusCompleted, Result: b}))
			},
		},
		{
			name: "should set error if result is not serializable",
			job:  jobs.Job{ID: "1"},
			workerExpectation: func(w *MockWorker) {
				w.EXPECT().Execute(gomock.Any(), gomock.Eq(jobs.Job{ID: "1"})).Return(map[string]interface{}{"foo": make(chan int)}, nil)
			},
			updaterExpectation: func(u *db.MockJobsStore) {
				u.EXPECT().Update(gomock.Any(), gomock.Eq(jobs.Job{ID: "1", Status: jobs.StatusCompleted, LastError: "json: unsupported type: chan int"}))
			},
		},
		{
			name: "should retry if job fails and retry is available",
			job:  jobs.Job{ID: "1", Options: jobs.Options{RetryInterval: time.Second, MaxRetries: 1}},
			workerExpectation: func(w *MockWorker) {
				w.EXPECT().Execute(gomock.Any(), gomock.Eq(jobs.Job{
					ID:      "1",
					Options: jobs.Options{RetryInterval: time.Second, MaxRetries: 1},
				})).Return(nil, errors.New("exec error"))
			},
			updaterExpectation: func(u *db.MockJobsStore) {
				now = func() time.Time {
					return time.Time{}
				}
				u.EXPECT().Update(gomock.Any(), gomock.Eq(jobs.Job{
					ID:          "1",
					Options:     jobs.Options{RetryInterval: time.Second, MaxRetries: 1},
					RetryCount:  1,
					Status:      jobs.StatusScheduled,
					LastError:   "exec error",
					ScheduledAt: time.Time{}.Add(time.Second),
				}))
			},
		},
		{
			name: "should fail if job fails and retry is not available",
			job:  jobs.Job{ID: "1", RetryCount: 1, Options: jobs.Options{RetryInterval: time.Second, MaxRetries: 1}},
			workerExpectation: func(w *MockWorker) {
				w.EXPECT().Execute(gomock.Any(), gomock.Eq(jobs.Job{
					ID:         "1",
					RetryCount: 1,
					Options:    jobs.Options{RetryInterval: time.Second, MaxRetries: 1},
				})).Return(nil, errors.New("exec error"))
			},
			updaterExpectation: func(u *db.MockJobsStore) {
				now = func() time.Time {
					return time.Time{}
				}
				u.EXPECT().Update(gomock.Any(), gomock.Eq(jobs.Job{
					ID:         "1",
					Options:    jobs.Options{RetryInterval: time.Second, MaxRetries: 1},
					RetryCount: 1,
					Status:     jobs.StatusFailed,
					LastError:  "exec error",
				}))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			w := NewMockWorker(mockCtrl)
			tt.workerExpectation(w)

			mockupdaterCtrl := gomock.NewController(t)
			defer mockupdaterCtrl.Finish()
			store := db.NewMockJobsStore(mockupdaterCtrl)
			tt.updaterExpectation(store)

			r := &handler{
				worker:  w,
				updater: store,
			}
			err := r.Handle(context.TODO(), tt.job)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
