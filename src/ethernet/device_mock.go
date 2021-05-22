// Code generated by MockGen. DO NOT EDIT.
// Source: device.go

// Package ethernet is a generated GoMock package.
package ethernet

import (
	reflect "reflect"

	error "github.com/42milez/ProtocolStack/src/error"
	gomock "github.com/golang/mock/gomock"
)

// MockIDevice is a mock of IDevice interface.
type MockIDevice struct {
	ctrl     *gomock.Controller
	recorder *MockIDeviceMockRecorder
}

// MockIDeviceMockRecorder is the mock recorder for MockIDevice.
type MockIDeviceMockRecorder struct {
	mock *MockIDevice
}

// NewMockIDevice creates a new mock instance.
func NewMockIDevice(ctrl *gomock.Controller) *MockIDevice {
	mock := &MockIDevice{ctrl: ctrl}
	mock.recorder = &MockIDeviceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDevice) EXPECT() *MockIDeviceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockIDevice) Close() error.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error.Error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockIDeviceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIDevice)(nil).Close))
}

// Disable mocks base method.
func (m *MockIDevice) Disable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Disable")
}

// Disable indicates an expected call of Disable.
func (mr *MockIDeviceMockRecorder) Disable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disable", reflect.TypeOf((*MockIDevice)(nil).Disable))
}

// Enable mocks base method.
func (m *MockIDevice) Enable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Enable")
}

// Enable indicates an expected call of Enable.
func (mr *MockIDeviceMockRecorder) Enable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enable", reflect.TypeOf((*MockIDevice)(nil).Enable))
}

// Equal mocks base method.
func (m *MockIDevice) Equal(dev IDevice) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Equal", dev)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Equal indicates an expected call of Equal.
func (mr *MockIDeviceMockRecorder) Equal(dev interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Equal", reflect.TypeOf((*MockIDevice)(nil).Equal), dev)
}

// Info mocks base method.
func (m *MockIDevice) Info() (string, string, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	return ret0, ret1, ret2
}

// Info indicates an expected call of Info.
func (mr *MockIDeviceMockRecorder) Info() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockIDevice)(nil).Info))
}

// IsUp mocks base method.
func (m *MockIDevice) IsUp() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUp")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsUp indicates an expected call of IsUp.
func (mr *MockIDeviceMockRecorder) IsUp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUp", reflect.TypeOf((*MockIDevice)(nil).IsUp))
}

// Open mocks base method.
func (m *MockIDevice) Open() error.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(error.Error)
	return ret0
}

// Open indicates an expected call of Open.
func (mr *MockIDeviceMockRecorder) Open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockIDevice)(nil).Open))
}

// Poll mocks base method.
func (m *MockIDevice) Poll(terminate bool) error.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poll", terminate)
	ret0, _ := ret[0].(error.Error)
	return ret0
}

// Poll indicates an expected call of Poll.
func (mr *MockIDeviceMockRecorder) Poll(terminate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poll", reflect.TypeOf((*MockIDevice)(nil).Poll), terminate)
}

// Transmit mocks base method.
func (m *MockIDevice) Transmit() error.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transmit")
	ret0, _ := ret[0].(error.Error)
	return ret0
}

// Transmit indicates an expected call of Transmit.
func (mr *MockIDeviceMockRecorder) Transmit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transmit", reflect.TypeOf((*MockIDevice)(nil).Transmit))
}