//go:generate ${LOCAL_BIN}/mockgen -source ./transactor.go -destination=./mocks/transactor.go -package=mock_transactor
package transactor

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/khostya/pvz/pkg/postgres/exec"
)

const key = "tx"

type (
	// Transactor .
	Transactor interface {
		RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
		RunReadCommited(ctx context.Context, fx func(ctxTX context.Context) error) error
		Unwrap(err error) error
		GetQueryEngine(ctx context.Context) exec.QueryEngine
	}

	QueryEngineProvider interface {
		GetQueryEngine(ctx context.Context) exec.QueryEngine
	}

	Pool interface {
		exec.QueryEngine
		BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	}

	TransactionManager struct {
		pool Pool
	}
)

func NewTransactionManager(pool Pool) *TransactionManager {
	return &TransactionManager{pool}
}

func (tm *TransactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	options := pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	}

	return tm.run(ctx, fx, options)
}

func (tm *TransactionManager) RunReadCommited(ctx context.Context, fx func(ctxTX context.Context) error) error {
	options := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	return tm.run(ctx, fx, options)
}

func (tm *TransactionManager) run(ctx context.Context, fx func(ctx context.Context) error, options pgx.TxOptions) error {
	tx, err := tm.pool.BeginTx(ctx, options)

	if err != nil {
		return TransactionError{Inner: err}
	}
	if err := fx(context.WithValue(ctx, key, tx)); err != nil {
		return TransactionError{Inner: err, Rollback: tx.Rollback(ctx)}
	}

	if err := tx.Commit(ctx); err != nil {
		return TransactionError{Inner: err, Rollback: tx.Rollback(ctx)}
	}

	return nil
}

func (tm *TransactionManager) Unwrap(err error) error {
	if err == nil {
		return nil
	}

	var transactionError TransactionError
	ok := errors.As(err, &transactionError)
	if !ok {
		return err
	}
	return transactionError.Inner
}

func (tm *TransactionManager) GetQueryEngine(ctx context.Context) exec.QueryEngine {
	tx, ok := ctx.Value(key).(exec.QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.pool
}
