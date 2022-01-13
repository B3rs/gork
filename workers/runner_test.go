package workers

import (
	"context"
	"testing"

	"github.com/B3rs/gork/jobs"
	gomock "github.com/golang/mock/gomock"
)

func Test_runner_Run(t *testing.T) {
	type fields struct {
		worker Worker
		update updateJobFunc
	}
	type args struct {
		ctx context.Context
		job *jobs.Job
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			w := NewMockWorker(mockCtrl)

			r := &runner{
				worker: tt.fields.worker,
				update: tt.fields.update,
			}
			if err := r.Run(context.TODO(), tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("runner.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
