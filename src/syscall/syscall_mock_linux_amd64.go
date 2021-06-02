// Code generated by MockGen. DO NOT EDIT.
// Source: syscall_linux_amd64.go

// Package syscall is a generated GoMock package.
package syscall

import (
	reflect "reflect"
	syscall "syscall"
	unsafe "unsafe"

	gomock "github.com/golang/mock/gomock"
)

// MockISyscall is a mock of ISyscall interface.
type MockISyscall struct {
	ctrl     *gomock.Controller
	recorder *MockISyscallMockRecorder
}

// MockISyscallMockRecorder is the mock recorder for MockISyscall.
type MockISyscallMockRecorder struct {
	mock *MockISyscall
}

// NewMockISyscall creates a new mock instance.
func NewMockISyscall(ctrl *gomock.Controller) *MockISyscall {
	mock := &MockISyscall{ctrl: ctrl}
	mock.recorder = &MockISyscallMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISyscall) EXPECT() *MockISyscallMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockISyscall) Close(fd int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", fd)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockISyscallMockRecorder) Close(fd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockISyscall)(nil).Close), fd)
}

// EpollCreate1 mocks base method.
func (m *MockISyscall) EpollCreate1(flag int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EpollCreate1", flag)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EpollCreate1 indicates an expected call of EpollCreate1.
func (mr *MockISyscallMockRecorder) EpollCreate1(flag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EpollCreate1", reflect.TypeOf((*MockISyscall)(nil).EpollCreate1), flag)
}

// EpollCtl mocks base method.
func (m *MockISyscall) EpollCtl(epfd, op, fd int, event *syscall.EpollEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EpollCtl", epfd, op, fd, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// EpollCtl indicates an expected call of EpollCtl.
func (mr *MockISyscallMockRecorder) EpollCtl(epfd, op, fd, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EpollCtl", reflect.TypeOf((*MockISyscall)(nil).EpollCtl), epfd, op, fd, event)
}

// EpollWait mocks base method.
func (m *MockISyscall) EpollWait(epfd int, events []syscall.EpollEvent, msec int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EpollWait", epfd, events, msec)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EpollWait indicates an expected call of EpollWait.
func (mr *MockISyscallMockRecorder) EpollWait(epfd, events, msec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EpollWait", reflect.TypeOf((*MockISyscall)(nil).EpollWait), epfd, events, msec)
}

// Ioctl mocks base method.
func (m *MockISyscall) Ioctl(fd, code int, data unsafe.Pointer) syscall.Errno {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ioctl", fd, code, data)
	ret0, _ := ret[0].(syscall.Errno)
	return ret0
}

// Ioctl indicates an expected call of Ioctl.
func (mr *MockISyscallMockRecorder) Ioctl(fd, code, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ioctl", reflect.TypeOf((*MockISyscall)(nil).Ioctl), fd, code, data)
}

// Open mocks base method.
func (m *MockISyscall) Open(path string, mode int, perm uint32) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", path, mode, perm)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *MockISyscallMockRecorder) Open(path, mode, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockISyscall)(nil).Open), path, mode, perm)
}

// Read mocks base method.
func (m *MockISyscall) Read(fd int, p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", fd, p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockISyscallMockRecorder) Read(fd, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockISyscall)(nil).Read), fd, p)
}

// Socket mocks base method.
func (m *MockISyscall) Socket(domain, typ, proto int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Socket", domain, typ, proto)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Socket indicates an expected call of Socket.
func (mr *MockISyscallMockRecorder) Socket(domain, typ, proto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Socket", reflect.TypeOf((*MockISyscall)(nil).Socket), domain, typ, proto)
}

// Write mocks base method.
func (m *MockISyscall) Write(fd int, p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", fd, p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockISyscallMockRecorder) Write(fd, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockISyscall)(nil).Write), fd, p)
}
