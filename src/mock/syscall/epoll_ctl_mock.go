// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/k12n3ud0n/Workspace/ProtocolStack/src/syscall/epoll_ctl.go

// Package mock_syscall is a generated GoMock package.
package mock_syscall

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEpollCtlSyscallInterface is a mock of EpollCtlSyscallInterface interface.
type MockEpollCtlSyscallInterface struct {
	ctrl     *gomock.Controller
	recorder *MockEpollCtlSyscallInterfaceMockRecorder
}

// MockEpollCtlSyscallInterfaceMockRecorder is the mock recorder for MockEpollCtlSyscallInterface.
type MockEpollCtlSyscallInterfaceMockRecorder struct {
	mock *MockEpollCtlSyscallInterface
}

// NewMockEpollCtlSyscallInterface creates a new mock instance.
func NewMockEpollCtlSyscallInterface(ctrl *gomock.Controller) *MockEpollCtlSyscallInterface {
	mock := &MockEpollCtlSyscallInterface{ctrl: ctrl}
	mock.recorder = &MockEpollCtlSyscallInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEpollCtlSyscallInterface) EXPECT() *MockEpollCtlSyscallInterfaceMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockEpollCtlSyscallInterface) Exec() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec")
	ret0, _ := ret[0].(error)
	return ret0
}

// Exec indicates an expected call of Exec.
func (mr *MockEpollCtlSyscallInterfaceMockRecorder) Exec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockEpollCtlSyscallInterface)(nil).Exec))
}
