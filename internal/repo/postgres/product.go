//go:generate ${LOCAL_BIN}/ifacemaker -f ./product.go -s ProductRepo -i ProductRepo -p mock_postgres -o ./mocks/product.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/product.go -destination=./mocks/product.go -package=mock_postgres

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

type ProductRepo struct {
	db transactor.QueryEngineProvider
}

func NewProductRepo(db transactor.QueryEngineProvider) *ProductRepo {
	return &ProductRepo{db}
}

func (r ProductRepo) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	db := r.db.GetQueryEngine(ctx)

	record := schema.NewProduct(product)

	query := sq.Insert(schema.Product{}.TableName()).
		Columns(record.InsertColumns()...).
		Values(record.Values()...).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(record.Columns(), ", ")))

	res, err := exec.ScanOne[schema.Product](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainProduct(&res), nil
}

func (r ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	db := r.db.GetQueryEngine(ctx)

	query := sq.Select(schema.Product{}.Columns()...).
		From(schema.Product{}.TableName()).
		Where("id = $1", id).
		PlaceholderFormat(sq.Dollar)

	res, err := exec.ScanOne[*schema.Product](ctx, query, db)
	if err != nil {
		return nil, err
	}

	return schema.NewDomainProduct(res), nil
}

func (r ProductRepo) DeleteLastByDateTimeAndReceptionID(ctx context.Context, receptionID uuid.UUID) error {
	db := r.db.GetQueryEngine(ctx)

	lastProductByDateTime := sq.Select("p.id").
		From(schema.Product{}.TableName()+" AS p").
		Where("p.reception_id = ?", receptionID).
		OrderBy("p.date_time desc").
		Limit(1)

	deleteQuery := sq.Delete(schema.Product{}.TableName()).
		Where(sq.Expr("id = (?)", lastProductByDateTime)).
		PlaceholderFormat(sq.Dollar)

	err := exec.Delete(ctx, deleteQuery, db)
	return err
}
