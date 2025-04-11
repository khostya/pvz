package schema

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

type (
	Product struct {
		ID          uuid.UUID `db:"products.id"`
		Type        string    `db:"products.type"`
		ReceptionID uuid.UUID `db:"products.reception_id"`
		DateTime    time.Time `db:"products.date_time"`
	}

	NullableProduct struct {
		ID          *uuid.UUID `db:"products.id"`
		Type        *string    `db:"products.type"`
		ReceptionID *uuid.UUID `db:"products.reception_id"`
		DateTime    *time.Time `db:"products.date_time"`
	}
)

func NewProduct(d *domain.Product) *Product {
	return &Product{
		ID:          d.ID,
		Type:        string(d.Type),
		ReceptionID: d.ReceptionID,
	}
}

func NewDomainProduct(d *Product) *domain.Product {
	return &domain.Product{
		ID:          d.ID,
		DateTime:    d.DateTime,
		Type:        domain.ProductType(d.Type),
		ReceptionID: d.ReceptionID,
	}
}

func (Product) TableName() string {
	return "products"
}

func (p Product) InsertColumns() []string {
	return []string{"id", "type", "reception_id", "date_time"}
}

func (p Product) Columns() []string {
	return []string{"products.id as \"products.id\"", "products.type as \"products.type\"",
		"products.reception_id as \"products.reception_id\"",
		"products.date_time as \"products.date_time\""}
}

func (p Product) Values() []any {
	return []any{p.ID, p.Type, p.ReceptionID, p.DateTime}
}
