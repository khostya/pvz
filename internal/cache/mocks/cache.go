// Code generated by MockGen. DO NOT EDIT.
// Source: ./cache.go
//
// Generated by this command:
//
//	mockgen -source ./cache.go -destination=./mocks/cache.go -package=mock_cache
//

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCache is a mock of GetPVZListResponse interface.
type MockCache[K comparable, V any] struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder[K, V]
	isgomock struct{}
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder[K comparable, V any] struct {
	mock *MockCache[K, V]
}

// NewMockCache creates a new mock instance.
func NewMockCache[K comparable, V any](ctrl *gomock.Controller) *MockCache[K, V] {
	mock := &MockCache[K, V]{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder[K, V]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache[K, V]) EXPECT() *MockCacheMockRecorder[K, V] {
	return m.recorder
}

// Get mocks base method.
func (m *MockCache[K, V]) Get(arg0 K) (V, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(V)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder[K, V]) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache[K, V])(nil).Get), arg0)
}

// Put mocks base method.
func (m *MockCache[K, V]) Put(arg0 K, arg1 V) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Put", arg0, arg1)
}

// Put indicates an expected call of Put.
func (mr *MockCacheMockRecorder[K, V]) Put(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockCache[K, V])(nil).Put), arg0, arg1)
}
