package workers

import (
	context "context"
	"errors"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_reaper_Run(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error, 3)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	q := NewMockQueue(mockCtrl)

	first := q.EXPECT().RequeueTimedOutJobs(gomock.Any(), gomock.Eq(7*time.Second)).Return(nil).Times(1)
	q.EXPECT().RequeueTimedOutJobs(gomock.Any(), gomock.Eq(7*time.Second)).Return(errors.New("reaper error")).Times(1).After(first).Do(func(_, _ interface{}) { cancel() })

	r := newReaper(q, 1*time.Millisecond, 7*time.Second)
	r.Run(ctx, errChan)

	assert.Equal(t, errors.New("reaper error"), <-errChan)
}
