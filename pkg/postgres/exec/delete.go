package exec

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/khostya/pvz/pkg/postgres/repoerr"
	pgxtransactor "github.com/khostya/pvz/pkg/postgres/transactor"
)

func Delete(ctx context.Context, query sq.Sqlizer, db pgxtransactor.QueryEngine) error {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrNotFound
	}
	return nil
}
