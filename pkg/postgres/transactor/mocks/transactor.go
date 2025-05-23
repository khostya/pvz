// Code generated by MockGen. DO NOT EDIT.
// Source: ./transactor.go
//
// Generated by this command:
//
//	mockgen -source ./transactor.go -destination=./mocks/transactor.go -package=mock_transactor
//

// Package mock_transactor is a generated GoMock package.
package mock_transactor

import (
	context "context"
	reflect "reflect"

	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	exec "github.com/khostya/pvz/pkg/postgres/exec"
	gomock "go.uber.org/mock/gomock"
)

// MockTransactor is a mock of Transactor interface.
type MockTransactor struct {
	ctrl     *gomock.Controller
	recorder *MockTransactorMockRecorder
	isgomock struct{}
}

// MockTransactorMockRecorder is the mock recorder for MockTransactor.
type MockTransactorMockRecorder struct {
	mock *MockTransactor
}

// NewMockTransactor creates a new mock instance.
func NewMockTransactor(ctrl *gomock.Controller) *MockTransactor {
	mock := &MockTransactor{ctrl: ctrl}
	mock.recorder = &MockTransactorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactor) EXPECT() *MockTransactorMockRecorder {
	return m.recorder
}

// GetQueryEngine mocks base method.
func (m *MockTransactor) GetQueryEngine(ctx context.Context) exec.QueryEngine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueryEngine", ctx)
	ret0, _ := ret[0].(exec.QueryEngine)
	return ret0
}

// GetQueryEngine indicates an expected call of GetQueryEngine.
func (mr *MockTransactorMockRecorder) GetQueryEngine(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueryEngine", reflect.TypeOf((*MockTransactor)(nil).GetQueryEngine), ctx)
}

// RunReadCommited mocks base method.
func (m *MockTransactor) RunReadCommited(ctx context.Context, fx func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunReadCommited", ctx, fx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunReadCommited indicates an expected call of RunReadCommited.
func (mr *MockTransactorMockRecorder) RunReadCommited(ctx, fx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunReadCommited", reflect.TypeOf((*MockTransactor)(nil).RunReadCommited), ctx, fx)
}

// RunRepeatableRead mocks base method.
func (m *MockTransactor) RunRepeatableRead(ctx context.Context, fx func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunRepeatableRead", ctx, fx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunRepeatableRead indicates an expected call of RunRepeatableRead.
func (mr *MockTransactorMockRecorder) RunRepeatableRead(ctx, fx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunRepeatableRead", reflect.TypeOf((*MockTransactor)(nil).RunRepeatableRead), ctx, fx)
}

// Unwrap mocks base method.
func (m *MockTransactor) Unwrap(err error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unwrap", err)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unwrap indicates an expected call of Unwrap.
func (mr *MockTransactorMockRecorder) Unwrap(err any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unwrap", reflect.TypeOf((*MockTransactor)(nil).Unwrap), err)
}

// MockQueryEngineProvider is a mock of QueryEngineProvider interface.
type MockQueryEngineProvider struct {
	ctrl     *gomock.Controller
	recorder *MockQueryEngineProviderMockRecorder
	isgomock struct{}
}

// MockQueryEngineProviderMockRecorder is the mock recorder for MockQueryEngineProvider.
type MockQueryEngineProviderMockRecorder struct {
	mock *MockQueryEngineProvider
}

// NewMockQueryEngineProvider creates a new mock instance.
func NewMockQueryEngineProvider(ctrl *gomock.Controller) *MockQueryEngineProvider {
	mock := &MockQueryEngineProvider{ctrl: ctrl}
	mock.recorder = &MockQueryEngineProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueryEngineProvider) EXPECT() *MockQueryEngineProviderMockRecorder {
	return m.recorder
}

// GetQueryEngine mocks base method.
func (m *MockQueryEngineProvider) GetQueryEngine(ctx context.Context) exec.QueryEngine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueryEngine", ctx)
	ret0, _ := ret[0].(exec.QueryEngine)
	return ret0
}

// GetQueryEngine indicates an expected call of GetQueryEngine.
func (mr *MockQueryEngineProviderMockRecorder) GetQueryEngine(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueryEngine", reflect.TypeOf((*MockQueryEngineProvider)(nil).GetQueryEngine), ctx)
}

// MockPool is a mock of Pool interface.
type MockPool struct {
	ctrl     *gomock.Controller
	recorder *MockPoolMockRecorder
	isgomock struct{}
}

// MockPoolMockRecorder is the mock recorder for MockPool.
type MockPoolMockRecorder struct {
	mock *MockPool
}

// NewMockPool creates a new mock instance.
func NewMockPool(ctrl *gomock.Controller) *MockPool {
	mock := &MockPool{ctrl: ctrl}
	mock.recorder = &MockPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPool) EXPECT() *MockPoolMockRecorder {
	return m.recorder
}

// BeginTx mocks base method.
func (m *MockPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, txOptions)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockPoolMockRecorder) BeginTx(ctx, txOptions any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockPool)(nil).BeginTx), ctx, txOptions)
}

// Exec mocks base method.
func (m *MockPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockPoolMockRecorder) Exec(ctx, sql any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockPool)(nil).Exec), varargs...)
}

// Query mocks base method.
func (m *MockPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(pgx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockPoolMockRecorder) Query(ctx, sql any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockPool)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	m.ctrl.T.Helper()
	varargs := []any{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(pgx.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockPoolMockRecorder) QueryRow(ctx, sql any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockPool)(nil).QueryRow), varargs...)
}
