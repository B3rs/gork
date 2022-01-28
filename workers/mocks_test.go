// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/B3rs/gork/workers (interfaces: Queue,Worker,Handler,Spawner)

// Package workers is a generated GoMock package.
package workers

import (
	context "context"
	reflect "reflect"
	time "time"

	jobs "github.com/B3rs/gork/jobs"
	gomock "github.com/golang/mock/gomock"
)

// MockQueue is a mock of Queue interface.
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue.
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance.
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// Poll mocks base method.
func (m *MockQueue) Poll(arg0 context.Context) (jobs.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poll", arg0)
	ret0, _ := ret[0].(jobs.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Poll indicates an expected call of Poll.
func (mr *MockQueueMockRecorder) Poll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poll", reflect.TypeOf((*MockQueue)(nil).Poll), arg0)
}

// RequeueTimedOutJobs mocks base method.
func (m *MockQueue) RequeueTimedOutJobs(arg0 context.Context, arg1 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequeueTimedOutJobs", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RequeueTimedOutJobs indicates an expected call of RequeueTimedOutJobs.
func (mr *MockQueueMockRecorder) RequeueTimedOutJobs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequeueTimedOutJobs", reflect.TypeOf((*MockQueue)(nil).RequeueTimedOutJobs), arg0, arg1)
}

// MockWorker is a mock of Worker interface.
type MockWorker struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerMockRecorder
}

// MockWorkerMockRecorder is the mock recorder for MockWorker.
type MockWorkerMockRecorder struct {
	mock *MockWorker
}

// NewMockWorker creates a new mock instance.
func NewMockWorker(ctrl *gomock.Controller) *MockWorker {
	mock := &MockWorker{ctrl: ctrl}
	mock.recorder = &MockWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorker) EXPECT() *MockWorkerMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockWorker) Execute(arg0 context.Context, arg1 jobs.Job) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockWorkerMockRecorder) Execute(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockWorker)(nil).Execute), arg0, arg1)
}

// OnFailure mocks base method.
func (m *MockWorker) OnFailure(arg0 context.Context, arg1 jobs.Job) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnFailure", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnFailure indicates an expected call of OnFailure.
func (mr *MockWorkerMockRecorder) OnFailure(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnFailure", reflect.TypeOf((*MockWorker)(nil).OnFailure), arg0, arg1)
}

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockHandler) Handle(arg0 context.Context, arg1 jobs.Job) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockHandlerMockRecorder) Handle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockHandler)(nil).Handle), arg0, arg1)
}

// MockSpawner is a mock of Spawner interface.
type MockSpawner struct {
	ctrl     *gomock.Controller
	recorder *MockSpawnerMockRecorder
}

// MockSpawnerMockRecorder is the mock recorder for MockSpawner.
type MockSpawnerMockRecorder struct {
	mock *MockSpawner
}

// NewMockSpawner creates a new mock instance.
func NewMockSpawner(ctrl *gomock.Controller) *MockSpawner {
	mock := &MockSpawner{ctrl: ctrl}
	mock.recorder = &MockSpawnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpawner) EXPECT() *MockSpawnerMockRecorder {
	return m.recorder
}

// Done mocks base method.
func (m *MockSpawner) Done() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Done")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Done indicates an expected call of Done.
func (mr *MockSpawnerMockRecorder) Done() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockSpawner)(nil).Done))
}

// Shutdown mocks base method.
func (m *MockSpawner) Shutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockSpawnerMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockSpawner)(nil).Shutdown))
}

// Spawn mocks base method.
func (m *MockSpawner) Spawn(arg0 runner) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Spawn", arg0)
}

// Spawn indicates an expected call of Spawn.
func (mr *MockSpawnerMockRecorder) Spawn(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Spawn", reflect.TypeOf((*MockSpawner)(nil).Spawn), arg0)
}

// Wait mocks base method.
func (m *MockSpawner) Wait() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Wait")
}

// Wait indicates an expected call of Wait.
func (mr *MockSpawnerMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockSpawner)(nil).Wait))
}
