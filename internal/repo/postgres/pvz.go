//go:generate ${LOCAL_BIN}/ifacemaker -f ./pvz.go -s PvzRepo -i PvzRepo -p mock_postgres -o ./mocks/pvz.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/pvz.go -destination=./mocks/pvz.go -package=mock_postgres

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

type PvzRepo struct {
	db transactor.QueryEngineProvider
}

func NewPvzRepo(db transactor.QueryEngineProvider) *PvzRepo {
	return &PvzRepo{db}
}

func (r PvzRepo) Create(ctx context.Context, pvz *domain.PVZ) (*domain.PVZ, error) {
	db := r.db.GetQueryEngine(ctx)

	record := schema.NewPVZ(pvz)

	query := sq.Insert(schema.PVZ{}.TableName()).
		Columns(record.Columns()...).
		Values(record.Values()...).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(record.Columns(), ", ")))

	res, err := exec.ScanOne[schema.PVZ](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainPVZ(&res), nil
}

func (r PvzRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.PVZ, error) {
	db := r.db.GetQueryEngine(ctx)

	query := sq.Select(schema.PVZ{}.Columns()...).
		From(schema.PVZ{}.TableName()).
		Where("id = $1", id).
		PlaceholderFormat(sq.Dollar)

	res, err := exec.ScanOne[schema.PVZ](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainPVZ(&res), nil
}
