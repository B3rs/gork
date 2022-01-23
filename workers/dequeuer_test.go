package workers

import (
	"context"
	"errors"
	"testing"

	jobs "github.com/B3rs/gork/jobs"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Dequeuer_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errsChan := make(chan error, 3)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	q := NewMockQueue(mockCtrl)

	runnerMockCtrl := gomock.NewController(t)
	defer runnerMockCtrl.Finish()
	r := NewMockHandler(runnerMockCtrl)

	dequeueCall := q.EXPECT().Dequeue(gomock.Any()).Return(jobs.Job{ID: "job1"}, nil).Times(1)
	handleCall := r.EXPECT().Handle(gomock.Any(), jobs.Job{ID: "job1"}).Return(nil).Times(1).After(dequeueCall)

	dequeueCall2 := q.EXPECT().Dequeue(gomock.Any()).Return(jobs.Job{ID: "job2"}, nil).Times(1).After(handleCall)
	handleCall2 := r.EXPECT().Handle(gomock.Any(), jobs.Job{ID: "job2"}).Return(nil).Times(1).After(dequeueCall2).After(handleCall)

	dequeueCall3 := q.EXPECT().Dequeue(gomock.Any()).Return(jobs.Job{ID: "job3"}, nil).Times(1).After(handleCall2).Return(jobs.Job{}, jobs.ErrJobNotFound)

	dequeueCall4 := q.EXPECT().Dequeue(gomock.Any()).Return(jobs.Job{}, errors.New("queue error")).Times(1).After(dequeueCall3)

	dequeueCall5 := q.EXPECT().Dequeue(gomock.Any()).Return(jobs.Job{ID: "job4"}, nil).Times(1).After(dequeueCall4)
	r.EXPECT().Handle(gomock.Any(), jobs.Job{ID: "job4"}).Return(errors.New("run error")).Times(1).After(dequeueCall5).Do(func(_, _ interface{}) { cancel() })

	s := &dequeuer{
		queue:   q,
		handler: r,
	}
	s.Run(ctx, errsChan)

	assert.Equal(t, errors.New("queue error"), <-errsChan)
	assert.Equal(t, errors.New("run error"), <-errsChan)
}
