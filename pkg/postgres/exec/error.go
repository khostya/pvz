package exec

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	if ok {
		return pgErr.Code == pgerrcode.UniqueViolation
	}
	return false
}
