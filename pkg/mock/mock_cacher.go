// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vatsal-chaturvedi/article-management-sys/internal/repo/cacher (interfaces: CacherI)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockCacherI is a mock of CacherI interface.
type MockCacherI struct {
	ctrl     *gomock.Controller
	recorder *MockCacherIMockRecorder
}

// MockCacherIMockRecorder is the mock recorder for MockCacherI.
type MockCacherIMockRecorder struct {
	mock *MockCacherI
}

// NewMockCacherI creates a new mock instance.
func NewMockCacherI(ctrl *gomock.Controller) *MockCacherI {
	mock := &MockCacherI{ctrl: ctrl}
	mock.recorder = &MockCacherIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacherI) EXPECT() *MockCacherIMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockCacherI) Get(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacherIMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCacherI)(nil).Get), arg0)
}

// Set mocks base method.
func (m *MockCacherI) Set(arg0 string, arg1 interface{}, arg2 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacherIMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCacherI)(nil).Set), arg0, arg1, arg2)
}
