// Code generated by MockGen. DO NOT EDIT.
// Source: ./mocks/product.go
//
// Generated by this command:
//
//	mockgen -source ./mocks/product.go -destination=./mocks/product.go -package=mock_postgres
//

// Package mock_postgres is a generated GoMock package.
package mock_postgres

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	domain "github.com/khostya/pvz/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockProductRepo is a mock of ProductRepo interface.
type MockProductRepo struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepoMockRecorder
	isgomock struct{}
}

// MockProductRepoMockRecorder is the mock recorder for MockProductRepo.
type MockProductRepoMockRecorder struct {
	mock *MockProductRepo
}

// NewMockProductRepo creates a new mock instance.
func NewMockProductRepo(ctrl *gomock.Controller) *MockProductRepo {
	mock := &MockProductRepo{ctrl: ctrl}
	mock.recorder = &MockProductRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepo) EXPECT() *MockProductRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProductRepo) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, product)
	ret0, _ := ret[0].(*domain.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockProductRepoMockRecorder) Create(ctx, product any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductRepo)(nil).Create), ctx, product)
}

// DeleteLastByDateTime mocks base method.
func (m *MockProductRepo) DeleteLastByDateTime(ctx context.Context, receptionID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastByDateTime", ctx, receptionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastByDateTime indicates an expected call of DeleteLastByDateTime.
func (mr *MockProductRepoMockRecorder) DeleteLastByDateTime(ctx, receptionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastByDateTime", reflect.TypeOf((*MockProductRepo)(nil).DeleteLastByDateTime), ctx, receptionID)
}

// GetByID mocks base method.
func (m *MockProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockProductRepoMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockProductRepo)(nil).GetByID), ctx, id)
}
