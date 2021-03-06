// Code generated by MockGen. DO NOT EDIT.
// Source: spawner.go

// Package workers is a generated GoMock package.
package workers

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Mockrunner is a mock of runner interface.
type Mockrunner struct {
	ctrl     *gomock.Controller
	recorder *MockrunnerMockRecorder
}

// MockrunnerMockRecorder is the mock recorder for Mockrunner.
type MockrunnerMockRecorder struct {
	mock *Mockrunner
}

// NewMockrunner creates a new mock instance.
func NewMockrunner(ctrl *gomock.Controller) *Mockrunner {
	mock := &Mockrunner{ctrl: ctrl}
	mock.recorder = &MockrunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrunner) EXPECT() *MockrunnerMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *Mockrunner) Run(arg0 context.Context, arg1 chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0, arg1)
}

// Run indicates an expected call of Run.
func (mr *MockrunnerMockRecorder) Run(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*Mockrunner)(nil).Run), arg0, arg1)
}
