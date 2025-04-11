// Code generated by MockGen. DO NOT EDIT.
// Source: ./mocks/pvz.go
//
// Generated by this command:
//
//	mockgen -source ./mocks/pvz.go -destination=./mocks/pvz.go -package=mock_postgres
//

// Package mock_postgres is a generated GoMock package.
package mock_postgres

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	domain "github.com/khostya/pvz/internal/domain"
	dto "github.com/khostya/pvz/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockPvzRepo is a mock of PvzRepo interface.
type MockPvzRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPvzRepoMockRecorder
	isgomock struct{}
}

// MockPvzRepoMockRecorder is the mock recorder for MockPvzRepo.
type MockPvzRepoMockRecorder struct {
	mock *MockPvzRepo
}

// NewMockPvzRepo creates a new mock instance.
func NewMockPvzRepo(ctrl *gomock.Controller) *MockPvzRepo {
	mock := &MockPvzRepo{ctrl: ctrl}
	mock.recorder = &MockPvzRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPvzRepo) EXPECT() *MockPvzRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPvzRepo) Create(ctx context.Context, pvz *domain.PVZ) (*domain.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, pvz)
	ret0, _ := ret[0].(*domain.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPvzRepoMockRecorder) Create(ctx, pvz any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPvzRepo)(nil).Create), ctx, pvz)
}

// GetAllPVZList mocks base method.
func (m *MockPvzRepo) GetAllPVZList(ctx context.Context) ([]*domain.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPVZList", ctx)
	ret0, _ := ret[0].([]*domain.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPVZList indicates an expected call of GetAllPVZList.
func (mr *MockPvzRepoMockRecorder) GetAllPVZList(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPVZList", reflect.TypeOf((*MockPvzRepo)(nil).GetAllPVZList), ctx)
}

// GetByID mocks base method.
func (m *MockPvzRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockPvzRepoMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockPvzRepo)(nil).GetByID), ctx, id)
}

// GetPVZ mocks base method.
func (m *MockPvzRepo) GetPVZ(ctx context.Context, param dto.GetPvzParam) ([]*domain.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPVZ", ctx, param)
	ret0, _ := ret[0].([]*domain.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPVZ indicates an expected call of GetPVZ.
func (mr *MockPvzRepoMockRecorder) GetPVZ(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPVZ", reflect.TypeOf((*MockPvzRepo)(nil).GetPVZ), ctx, param)
}
