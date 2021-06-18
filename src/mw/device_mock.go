// Code generated by MockGen. DO NOT EDIT.
// Source: device.go

// Package mw is a generated GoMock package.
package mw

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

// Addr mocks base method.
func (m *MockIDevice) Addr() EthAddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Addr")
	ret0, _ := ret[0].(EthAddr)
	return ret0
}

// Addr indicates an expected call of Addr.
func (mr *MockIDeviceMockRecorder) Addr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Addr", reflect.TypeOf((*MockIDevice)(nil).Addr))
}

// Close mocks base method.
func (m *MockIDevice) Close() error.E {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error.E)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockIDeviceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIDevice)(nil).Close))
}

// Down mocks base method.
func (m *MockIDevice) Down() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Down")
}

// Down indicates an expected call of Down.
func (mr *MockIDeviceMockRecorder) Down() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Down", reflect.TypeOf((*MockIDevice)(nil).Down))
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

// Flag mocks base method.
func (m *MockIDevice) Flag() DevFlag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "flag")
	ret0, _ := ret[0].(DevFlag)
	return ret0
}

// Flag indicates an expected call of Flag.
func (mr *MockIDeviceMockRecorder) Flag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "flag", reflect.TypeOf((*MockIDevice)(nil).Flag))
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

// MTU mocks base method.
func (m *MockIDevice) MTU() uint16 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MTU")
	ret0, _ := ret[0].(uint16)
	return ret0
}

// MTU indicates an expected call of MTU.
func (mr *MockIDeviceMockRecorder) MTU() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MTU", reflect.TypeOf((*MockIDevice)(nil).MTU))
}

// Name mocks base method.
func (m *MockIDevice) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockIDeviceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockIDevice)(nil).Name))
}

// Open mocks base method.
func (m *MockIDevice) Open() error.E {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(error.E)
	return ret0
}

// Open indicates an expected call of Open.
func (mr *MockIDeviceMockRecorder) Open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockIDevice)(nil).Open))
}

// Poll mocks base method.
func (m *MockIDevice) Poll() error.E {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poll")
	ret0, _ := ret[0].(error.E)
	return ret0
}

// Poll indicates an expected call of Poll.
func (mr *MockIDeviceMockRecorder) Poll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poll", reflect.TypeOf((*MockIDevice)(nil).Poll))
}

// Priv mocks base method.
func (m *MockIDevice) Priv() Privilege {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Priv")
	ret0, _ := ret[0].(Privilege)
	return ret0
}

// Priv indicates an expected call of Priv.
func (mr *MockIDeviceMockRecorder) Priv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Priv", reflect.TypeOf((*MockIDevice)(nil).Priv))
}

// Transmit mocks base method.
func (m *MockIDevice) Transmit(dst EthAddr, payload []byte, typ EthType) error.E {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transmit", dst, payload, typ)
	ret0, _ := ret[0].(error.E)
	return ret0
}

// Transmit indicates an expected call of Transmit.
func (mr *MockIDeviceMockRecorder) Transmit(dst, payload, typ interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transmit", reflect.TypeOf((*MockIDevice)(nil).Transmit), dst, payload, typ)
}

// Type mocks base method.
func (m *MockIDevice) Type() DevType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(DevType)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockIDeviceMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockIDevice)(nil).Type))
}

// Up mocks base method.
func (m *MockIDevice) Up() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Up")
}

// Up indicates an expected call of Up.
func (mr *MockIDeviceMockRecorder) Up() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Up", reflect.TypeOf((*MockIDevice)(nil).Up))
}
