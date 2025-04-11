package schema

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

type (
	Reception struct {
		schemaModel
		ID     uuid.UUID `db:"id"`
		PvzID  uuid.UUID `db:"pvz_id"`
		Status string    `db:"status"`
	}
)

func NewReception(d *domain.Reception) *Reception {
	return &Reception{
		schemaModel: schemaModel{
			CreatedAt: time.Now(),
		},
		ID:     d.ID,
		PvzID:  d.PvzId,
		Status: string(d.Status),
	}
}

func NewDomainReception(d *Reception) *domain.Reception {
	return &domain.Reception{
		ID:       d.ID,
		PvzId:    d.PvzID,
		DateTime: d.CreatedAt,
		Status:   domain.ReceptionStatus(d.Status),
	}
}

func (Reception) TableName() string {
	return "receptions"
}

func (p Reception) Columns() []string {
	return []string{"id", "pvz_id", "status", "created_at", "updated_at", "deleted_at"}
}

func (p Reception) Values() []any {
	return []any{p.ID, p.PvzID, p.Status, p.CreatedAt, p.UpdatedAt, p.DeletedAt}
}

func (p Reception) UpdateColumns() []string {
	return []string{"id", "pvz_id", "status", "updated_at"}
}

func (p Reception) UpdateValues() []any {
	return []any{p.ID, p.PvzID, p.Status, time.Now()}
}
