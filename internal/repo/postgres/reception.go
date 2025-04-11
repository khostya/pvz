//go:generate ${LOCAL_BIN}/ifacemaker -f ./reception.go -s ReceptionRepo -i ReceptionRepo -p mock_postgres -o ./mocks/reception.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/reception.go -destination=./mocks/reception.go -package=mock_postgres

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

type ReceptionRepo struct {
	db transactor.QueryEngineProvider
}

func NewReceptionRepo(db transactor.QueryEngineProvider) *ReceptionRepo {
	return &ReceptionRepo{db}
}

func (r ReceptionRepo) Create(ctx context.Context, reception *domain.Reception) (*domain.Reception, error) {
	db := r.db.GetQueryEngine(ctx)

	record := schema.NewReception(reception)

	query := sq.Insert(schema.Reception{}.TableName()).
		Columns(record.Columns()...).
		Values(record.Values()...).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(record.Columns(), ", ")))

	res, err := exec.ScanOne[schema.Reception](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainReception(&res), nil
}

func (r ReceptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reception, error) {
	return r.getByColumn(ctx, "id", id)
}

func (r ReceptionRepo) GetFirstByStatus(ctx context.Context, status domain.ReceptionStatus) (*domain.Reception, error) {
	return r.getByColumn(ctx, "status", status)
}

func (r ReceptionRepo) getByColumn(ctx context.Context, column string, value any) (*domain.Reception, error) {
	db := r.db.GetQueryEngine(ctx)

	query := sq.Select(schema.Reception{}.Columns()...).
		From(schema.Reception{}.TableName()).
		Where(fmt.Sprintf("%s = ?", column), value).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	res, err := exec.ScanOne[schema.Reception](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainReception(&res), nil
}

func (r ReceptionRepo) UpdateLastReceptionStatus(ctx context.Context, pvzID uuid.UUID, status domain.ReceptionStatus) (*domain.Reception, error) {
	db := r.db.GetQueryEngine(ctx)

	query := sq.Update(schema.Reception{}.TableName()).
		Where("pvz_id = ? and status != ?", pvzID, status).
		Set("status", status).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(schema.Reception{}.Columns(), ", ")))

	res, err := exec.ScanOne[schema.Reception](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainReception(&res), nil
}
