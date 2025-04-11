package transactor

import (
	"context"
	"errors"
	mock_query "github.com/khostya/pvz/pkg/postgres/exec/mocks"
	mock_transactor "github.com/khostya/pvz/pkg/postgres/transactor/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func newMock(t *testing.T) Pool {
	ctrl := gomock.NewController(t)
	return mock_transactor.NewMockPool(ctrl)
}

func TestTransactionManager_GetQueryEngine(t *testing.T) {
	ctrl := gomock.NewController(t)
	qe := mock_query.NewMockQueryEngine(ctrl)
	ctx := context.WithValue(context.Background(), key, qe)

	tm := TransactionManager{}
	actual := tm.GetQueryEngine(ctx)
	require.NotNil(t, actual)
}

func TestTransactionManager_Unwrap(t *testing.T) {
	m := newMock(t)
	tm := NewTransactionManager(m)

	t.Run("err is nil", func(t *testing.T) {
		err := tm.Unwrap(nil)
		require.NoError(t, err)
	})

	t.Run("err is not transaction error", func(t *testing.T) {
		expected := errors.New("err")
		actual := tm.Unwrap(expected)
		require.Equal(t, expected, actual)
	})

	t.Run("err is transaction error", func(t *testing.T) {
		innerErr := errors.New("inner")
		rollbackErr := errors.New("rollback")

		var err error = TransactionError{Inner: innerErr, Rollback: rollbackErr}
		actual := tm.Unwrap(err)

		require.Equal(t, innerErr, actual)
	})
}
