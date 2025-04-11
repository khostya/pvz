package schema

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

type (
	Product struct {
		schemaModel
		ID          uuid.UUID `db:"id"`
		Type        string    `db:"type"`
		PvzID       uuid.UUID `db:"pvz_id"`
		ReceptionID uuid.UUID `db:"reception_id"`
	}
)

func NewProduct(d *domain.Product) *Product {
	return &Product{
		schemaModel: schemaModel{
			CreatedAt: time.Now(),
		},
		ID:          d.ID,
		Type:        string(d.Type),
		PvzID:       d.PvzID,
		ReceptionID: d.ReceptionID,
	}
}

func NewDomainProduct(d *Product) *domain.Product {
	return &domain.Product{
		ID:          d.ID,
		DateTime:    d.CreatedAt,
		Type:        domain.ProductType(d.Type),
		PvzID:       d.PvzID,
		ReceptionID: d.ReceptionID,
	}
}

func (Product) TableName() string {
	return "products"
}

func (p Product) Columns() []string {
	return []string{"id", "type", "pvz_id", "reception_id", "created_at", "updated_at", "deleted_at"}
}

func (p Product) Values() []any {
	return []any{p.ID, p.Type, p.PvzID, p.ReceptionID, p.CreatedAt, p.UpdatedAt, p.DeletedAt}
}
