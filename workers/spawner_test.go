package workers

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func Test_spawner_Spawn(t *testing.T) {

	s := newSpawner(context.Background(), nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := NewMockrunner(ctrl)
	r.EXPECT().Run(gomock.Any(), gomock.Any()).Do(func(ctx, _ interface{}) {
		<-ctx.(context.Context).Done()
	})
	s.Spawn(r)

	r = NewMockrunner(ctrl)
	r.EXPECT().Run(gomock.Any(), gomock.Any()).Do(func(ctx, _ interface{}) {
		<-ctx.(context.Context).Done()
	})
	s.Spawn(r)

	r = NewMockrunner(ctrl)
	r.EXPECT().Run(gomock.Any(), gomock.Any()).Do(func(ctx, _ interface{}) {
		<-ctx.(context.Context).Done()
	})
	s.Spawn(r)

	s.Shutdown()

	s.Wait()
}
