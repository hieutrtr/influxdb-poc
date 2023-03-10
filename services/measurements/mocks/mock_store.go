// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hieutrtr/influxdb-poc/services/measurements/db (interfaces: Store)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	measurementdb "github.com/hieutrtr/influxdb-poc/services/measurements/db"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// ArchiveMeasurement mocks base method.
func (m *MockStore) ArchiveMeasurement(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveMeasurement", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ArchiveMeasurement indicates an expected call of ArchiveMeasurement.
func (mr *MockStoreMockRecorder) ArchiveMeasurement(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveMeasurement", reflect.TypeOf((*MockStore)(nil).ArchiveMeasurement), arg0, arg1)
}

// CreateMeasurement mocks base method.
func (m *MockStore) CreateMeasurement(arg0 context.Context, arg1 measurementdb.Measurement) (*measurementdb.MeasurementID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMeasurement", arg0, arg1)
	ret0, _ := ret[0].(*measurementdb.MeasurementID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMeasurement indicates an expected call of CreateMeasurement.
func (mr *MockStoreMockRecorder) CreateMeasurement(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMeasurement", reflect.TypeOf((*MockStore)(nil).CreateMeasurement), arg0, arg1)
}

// GetMeasurement mocks base method.
func (m *MockStore) GetMeasurement(arg0 context.Context, arg1 measurementdb.MeasurementID) (*measurementdb.Measurement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMeasurement", arg0, arg1)
	ret0, _ := ret[0].(*measurementdb.Measurement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMeasurement indicates an expected call of GetMeasurement.
func (mr *MockStoreMockRecorder) GetMeasurement(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMeasurement", reflect.TypeOf((*MockStore)(nil).GetMeasurement), arg0, arg1)
}

// ListMeasurements mocks base method.
func (m *MockStore) ListMeasurements(arg0 context.Context, arg1, arg2 int64) ([]measurementdb.Measurement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMeasurements", arg0, arg1, arg2)
	ret0, _ := ret[0].([]measurementdb.Measurement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMeasurements indicates an expected call of ListMeasurements.
func (mr *MockStoreMockRecorder) ListMeasurements(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMeasurements", reflect.TypeOf((*MockStore)(nil).ListMeasurements), arg0, arg1, arg2)
}
