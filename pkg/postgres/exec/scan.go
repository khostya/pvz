//go:generate ${LOCAL_BIN}/mockgen -source ./scan.go -destination=./mocks/query.go -package=mock_query
package exec

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/khostya/pvz/pkg/postgres/repoerr"
)

type (
	QueryEngine interface {
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	}
)

func ScanOne[T any](ctx context.Context, query sq.Sqlizer, db QueryEngine) (T, error) {
	var defaultT T

	records, err := ScanALL[T](ctx, query, db)
	if err != nil {
		return defaultT, err
	}

	return records[0], nil
}

func ScanALL[T any](ctx context.Context, query sq.Sqlizer, db QueryEngine) ([]T, error) {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var records []T
	err = pgxscan.Select(ctx, db, &records, rawQuery, args...)
	if IsDuplicateKeyError(err) {
		return nil, repoerr.ErrDuplicate
	}
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, repoerr.ErrNotFound
	}

	return records, nil
}
