package schema

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

type (
	Reception struct {
		ID       uuid.UUID `db:"receptions.id"`
		PvzID    uuid.UUID `db:"receptions.pvz_id"`
		Status   string    `db:"receptions.status"`
		DateTime time.Time `db:"receptions.date_time"`
	}
)

func NewReception(d *domain.Reception) *Reception {
	return &Reception{
		ID:       d.ID,
		PvzID:    d.PvzId,
		Status:   string(d.Status),
		DateTime: d.DateTime,
	}
}

func NewDomainReception(d *Reception) *domain.Reception {
	return &domain.Reception{
		ID:       d.ID,
		PvzId:    d.PvzID,
		DateTime: d.DateTime,
		Status:   domain.ReceptionStatus(d.Status),
	}
}

func (Reception) TableName() string {
	return "receptions"
}

func (p Reception) InsertColumns() []string {
	return []string{"id", "pvz_id", "status", "date_time"}
}

func (p Reception) Columns() []string {
	return []string{"receptions.id as \"receptions.id\"", "receptions.pvz_id as \"receptions.pvz_id\"",
		"receptions.status as \"receptions.status\"", "receptions.date_time as \"receptions.date_time\""}
}

func (p Reception) Values() []any {
	return []any{p.ID, p.PvzID, p.Status, p.DateTime}
}
