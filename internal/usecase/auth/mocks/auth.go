// Code generated by MockGen. DO NOT EDIT.
// Source: ./mocks/auth.go
//
// Generated by this command:
//
//	mockgen -source ./mocks/auth.go -destination=./mocks/auth.go -package=mock_auth
//

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	context "context"
	reflect "reflect"

	domain "github.com/khostya/pvz/internal/domain"
	dto "github.com/khostya/pvz/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
	isgomock struct{}
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// DummyLogin mocks base method.
func (m *MockAuth) DummyLogin(ctx context.Context, param dto.DummyLoginUserParam) (domain.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DummyLogin", ctx, param)
	ret0, _ := ret[0].(domain.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DummyLogin indicates an expected call of DummyLogin.
func (mr *MockAuthMockRecorder) DummyLogin(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DummyLogin", reflect.TypeOf((*MockAuth)(nil).DummyLogin), ctx, param)
}

// Login mocks base method.
func (m *MockAuth) Login(ctx context.Context, param dto.LoginUserParam) (domain.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, param)
	ret0, _ := ret[0].(domain.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthMockRecorder) Login(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuth)(nil).Login), ctx, param)
}

// Register mocks base method.
func (m *MockAuth) Register(ctx context.Context, param dto.RegisterUserParam) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, param)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockAuthMockRecorder) Register(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuth)(nil).Register), ctx, param)
}
