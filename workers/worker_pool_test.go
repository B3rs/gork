package workers

import (
	"testing"
)

func TestWorkerPool_Start(t *testing.T) {

	w := &WorkerPool{
		register:               newRegister(),
		db:                     tt.fields.db,
		errorHandler:           tt.fields.errorHandler,
		schedulerSleepInterval: tt.fields.schedulerSleepInterval,
		reaperInterval:         tt.fields.reaperInterval,
		shutdown:               tt.fields.shutdown,
	}
	w.Start()

}
