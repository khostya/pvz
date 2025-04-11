//go:generate ${LOCAL_BIN}/ifacemaker -f ./user.go -s UserRepo -i UserRepo -p mock_postgres -o ./mocks/user.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/user.go -destination=./mocks/user.go -package=mock_postgres

package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/repo/postgres/schema"
	"github.com/khostya/pvz/pkg/postgres/exec"
	"github.com/khostya/pvz/pkg/postgres/transactor"
	"strings"
)

type UserRepo struct {
	db transactor.QueryEngineProvider
}

func NewUserRepo(db transactor.QueryEngineProvider) *UserRepo {
	return &UserRepo{db}
}

func (r UserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	db := r.db.GetQueryEngine(ctx)

	record := schema.NewUser(user)

	query := sq.Insert(schema.User{}.TableName()).
		Columns(record.InsertColumns()...).
		Values(record.Values()...).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(record.Columns(), ", ")))

	res, err := exec.ScanOne[*schema.User](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainUser(res), nil
}

func (r UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return r.getByColumn(ctx, "id", id)
}

func (r UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.getByColumn(ctx, "email", email)
}

func (r UserRepo) getByColumn(ctx context.Context, column string, value any) (*domain.User, error) {
	db := r.db.GetQueryEngine(ctx)

	query := sq.Select(schema.User{}.Columns()...).
		From(schema.User{}.TableName()).
		Where(fmt.Sprintf("%s = ?", column), value).
		PlaceholderFormat(sq.Dollar)

	res, err := exec.ScanOne[schema.User](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainUser(&res), nil
}
