package exec

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/khostya/pvz/pkg/postgres/repoerr"
	"github.com/khostya/pvz/pkg/postgres/transactor"
)

func ScanOne[T any](ctx context.Context, query sq.Sqlizer, db transactor.QueryEngine) (T, error) {
	var defaultT T

	records, err := ScanALL[T](ctx, query, db)
	if IsDuplicateKeyError(err) {
		return defaultT, repoerr.ErrDuplicate
	}
	if err != nil {
		return defaultT, err
	}

	if len(records) == 0 {
		return defaultT, repoerr.ErrNotFound
	}

	return records[0], nil
}

func ScanALL[T any](ctx context.Context, query sq.Sqlizer, db transactor.QueryEngine) ([]T, error) {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var records []T
	return records, pgxscan.Select(ctx, db, &records, rawQuery, args...)
}
