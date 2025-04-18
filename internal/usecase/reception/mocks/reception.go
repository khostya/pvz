// Code generated by MockGen. DO NOT EDIT.
// Source: ./mocks/reception.go
//
// Generated by this command:
//
//	mockgen -source ./mocks/reception.go -destination=./mocks/reception.go -package=mock_reception
//

// Package mock_reception is a generated GoMock package.
package mock_reception

import (
	context "context"
	reflect "reflect"

	domain "github.com/khostya/pvz/internal/domain"
	dto "github.com/khostya/pvz/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockReception is a mock of Reception interface.
type MockReception struct {
	ctrl     *gomock.Controller
	recorder *MockReceptionMockRecorder
	isgomock struct{}
}

// MockReceptionMockRecorder is the mock recorder for MockReception.
type MockReceptionMockRecorder struct {
	mock *MockReception
}

// NewMockReception creates a new mock instance.
func NewMockReception(ctrl *gomock.Controller) *MockReception {
	mock := &MockReception{ctrl: ctrl}
	mock.recorder = &MockReceptionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReception) EXPECT() *MockReceptionMockRecorder {
	return m.recorder
}

// CloseLastReception mocks base method.
func (m *MockReception) CloseLastReception(ctx context.Context, param dto.CloseLastReceptionParam) (*domain.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseLastReception", ctx, param)
	ret0, _ := ret[0].(*domain.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseLastReception indicates an expected call of CloseLastReception.
func (mr *MockReceptionMockRecorder) CloseLastReception(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseLastReception", reflect.TypeOf((*MockReception)(nil).CloseLastReception), ctx, param)
}

// Create mocks base method.
func (m *MockReception) Create(ctx context.Context, param dto.CreateReceptionParam) (*domain.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, param)
	ret0, _ := ret[0].(*domain.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockReceptionMockRecorder) Create(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockReception)(nil).Create), ctx, param)
}

// DeleteLastProduct mocks base method.
func (m *MockReception) DeleteLastProduct(ctx context.Context, param dto.DeleteLastReceptionParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastProduct", ctx, param)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastProduct indicates an expected call of DeleteLastProduct.
func (mr *MockReceptionMockRecorder) DeleteLastProduct(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastProduct", reflect.TypeOf((*MockReception)(nil).DeleteLastProduct), ctx, param)
}
