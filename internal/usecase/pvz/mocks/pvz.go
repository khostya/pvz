// Code generated by MockGen. DO NOT EDIT.
// Source: ./mocks/pvz.go
//
// Generated by this command:
//
//	mockgen -source ./mocks/pvz.go -destination=./mocks/pvz.go -package=mock_pvz
//

// Package mock_pvz is a generated GoMock package.
package mock_pvz

import (
	context "context"
	reflect "reflect"

	domain "github.com/khostya/pvz/internal/domain"
	dto "github.com/khostya/pvz/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockPvz is a mock of Pvz interface.
type MockPvz struct {
	ctrl     *gomock.Controller
	recorder *MockPvzMockRecorder
	isgomock struct{}
}

// MockPvzMockRecorder is the mock recorder for MockPvz.
type MockPvzMockRecorder struct {
	mock *MockPvz
}

// NewMockPvz creates a new mock instance.
func NewMockPvz(ctrl *gomock.Controller) *MockPvz {
	mock := &MockPvz{ctrl: ctrl}
	mock.recorder = &MockPvzMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPvz) EXPECT() *MockPvzMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPvz) Create(ctx context.Context, param dto.CreatePvzParam) (*domain.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, param)
	ret0, _ := ret[0].(*domain.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPvzMockRecorder) Create(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPvz)(nil).Create), ctx, param)
}

// GetAllPvzList mocks base method.
func (m *MockPvz) GetAllPvzList(ctx context.Context) ([]*domain.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPvzList", ctx)
	ret0, _ := ret[0].([]*domain.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPvzList indicates an expected call of GetAllPvzList.
func (mr *MockPvzMockRecorder) GetAllPvzList(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPvzList", reflect.TypeOf((*MockPvz)(nil).GetAllPvzList), ctx)
}

// GetPvz mocks base method.
func (m *MockPvz) GetPvz(ctx context.Context, param dto.GetPvzParam) ([]*domain.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPvz", ctx, param)
	ret0, _ := ret[0].([]*domain.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPvz indicates an expected call of GetPvz.
func (mr *MockPvzMockRecorder) GetPvz(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPvz", reflect.TypeOf((*MockPvz)(nil).GetPvz), ctx, param)
}
