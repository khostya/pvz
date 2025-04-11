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
		Columns(record.InsertColumns()...).
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
	return r.getByColumn(ctx, sq.Eq{"id": id})
}

func (r ReceptionRepo) GetFirstByStatusAndPVZId(ctx context.Context, status domain.ReceptionStatus, pvzID uuid.UUID) (*domain.Reception, error) {
	return r.getByColumn(ctx, sq.Eq{"status": status, "pvz_id": pvzID})
}

func (r ReceptionRepo) getByColumn(ctx context.Context, eq sq.Eq) (*domain.Reception, error) {
	db := r.db.GetQueryEngine(ctx)

	query := sq.Select(schema.Reception{}.Columns()...).
		From(schema.Reception{}.TableName()).
		Where(eq).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	res, err := exec.ScanOne[schema.Reception](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainReception(&res), nil
}

func (r ReceptionRepo) UpdateReceptionStatusByID(ctx context.Context, id uuid.UUID, status domain.ReceptionStatus) (*domain.Reception, error) {
	db := r.db.GetQueryEngine(ctx)

	query := sq.Update(schema.Reception{}.TableName()).
		Where("id = ? and status != ?", id, status).
		Set("status", status).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(schema.Reception{}.Columns(), ", ")))

	res, err := exec.ScanOne[schema.Reception](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainReception(&res), nil
}
