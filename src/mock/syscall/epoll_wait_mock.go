// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/k12n3ud0n/Workspace/ProtocolStack/src/syscall/epoll_wait.go

// Package mock_syscall is a generated GoMock package.
package mock_syscall

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEpollWaitSyscallInterface is a mock of EpollWaitSyscallInterface interface.
type MockEpollWaitSyscallInterface struct {
	ctrl     *gomock.Controller
	recorder *MockEpollWaitSyscallInterfaceMockRecorder
}

// MockEpollWaitSyscallInterfaceMockRecorder is the mock recorder for MockEpollWaitSyscallInterface.
type MockEpollWaitSyscallInterfaceMockRecorder struct {
	mock *MockEpollWaitSyscallInterface
}

// NewMockEpollWaitSyscallInterface creates a new mock instance.
func NewMockEpollWaitSyscallInterface(ctrl *gomock.Controller) *MockEpollWaitSyscallInterface {
	mock := &MockEpollWaitSyscallInterface{ctrl: ctrl}
	mock.recorder = &MockEpollWaitSyscallInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEpollWaitSyscallInterface) EXPECT() *MockEpollWaitSyscallInterfaceMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockEpollWaitSyscallInterface) Exec() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockEpollWaitSyscallInterfaceMockRecorder) Exec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockEpollWaitSyscallInterface)(nil).Exec))
}
