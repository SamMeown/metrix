// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/storage.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	storage "github.com/SamMeown/metrix/internal/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockMetricsStorageGetter is a mock of MetricsStorageGetter interface.
type MockMetricsStorageGetter struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsStorageGetterMockRecorder
}

// MockMetricsStorageGetterMockRecorder is the mock recorder for MockMetricsStorageGetter.
type MockMetricsStorageGetterMockRecorder struct {
	mock *MockMetricsStorageGetter
}

// NewMockMetricsStorageGetter creates a new mock instance.
func NewMockMetricsStorageGetter(ctrl *gomock.Controller) *MockMetricsStorageGetter {
	mock := &MockMetricsStorageGetter{ctrl: ctrl}
	mock.recorder = &MockMetricsStorageGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsStorageGetter) EXPECT() *MockMetricsStorageGetterMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockMetricsStorageGetter) GetAll() (storage.MetricsStorageSnapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].(storage.MetricsStorageSnapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMetricsStorageGetterMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMetricsStorageGetter)(nil).GetAll))
}

// GetCounter mocks base method.
func (m *MockMetricsStorageGetter) GetCounter(name string) (*int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter", name)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockMetricsStorageGetterMockRecorder) GetCounter(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockMetricsStorageGetter)(nil).GetCounter), name)
}

// GetGauge mocks base method.
func (m *MockMetricsStorageGetter) GetGauge(name string) (*float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge", name)
	ret0, _ := ret[0].(*float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGauge indicates an expected call of GetGauge.
func (mr *MockMetricsStorageGetterMockRecorder) GetGauge(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockMetricsStorageGetter)(nil).GetGauge), name)
}

// MockMetricsStorage is a mock of MetricsStorage interface.
type MockMetricsStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsStorageMockRecorder
}

// MockMetricsStorageMockRecorder is the mock recorder for MockMetricsStorage.
type MockMetricsStorageMockRecorder struct {
	mock *MockMetricsStorage
}

// NewMockMetricsStorage creates a new mock instance.
func NewMockMetricsStorage(ctrl *gomock.Controller) *MockMetricsStorage {
	mock := &MockMetricsStorage{ctrl: ctrl}
	mock.recorder = &MockMetricsStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsStorage) EXPECT() *MockMetricsStorageMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockMetricsStorage) GetAll() (storage.MetricsStorageSnapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].(storage.MetricsStorageSnapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMetricsStorageMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMetricsStorage)(nil).GetAll))
}

// GetCounter mocks base method.
func (m *MockMetricsStorage) GetCounter(name string) (*int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter", name)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockMetricsStorageMockRecorder) GetCounter(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockMetricsStorage)(nil).GetCounter), name)
}

// GetGauge mocks base method.
func (m *MockMetricsStorage) GetGauge(name string) (*float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge", name)
	ret0, _ := ret[0].(*float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGauge indicates an expected call of GetGauge.
func (mr *MockMetricsStorageMockRecorder) GetGauge(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockMetricsStorage)(nil).GetGauge), name)
}

// SetCounter mocks base method.
func (m *MockMetricsStorage) SetCounter(name string, value int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCounter", name, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCounter indicates an expected call of SetCounter.
func (mr *MockMetricsStorageMockRecorder) SetCounter(name, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCounter", reflect.TypeOf((*MockMetricsStorage)(nil).SetCounter), name, value)
}

// SetGauge mocks base method.
func (m *MockMetricsStorage) SetGauge(name string, value float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetGauge", name, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetGauge indicates an expected call of SetGauge.
func (mr *MockMetricsStorageMockRecorder) SetGauge(name, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGauge", reflect.TypeOf((*MockMetricsStorage)(nil).SetGauge), name, value)
}
