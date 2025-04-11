package schema

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

type (
	PVZ struct {
		schemaModel
		ID               uuid.UUID `db:"id"`
		City             string    `db:"city"`
		RegistrationDate time.Time `db:"registration_date"`
	}
)

func NewPVZ(d *domain.PVZ) *PVZ {
	return &PVZ{
		schemaModel: schemaModel{
			CreatedAt: time.Now(),
		},
		ID:               d.ID,
		City:             string(d.City),
		RegistrationDate: d.RegistrationDate,
	}
}

func NewDomainPVZ(d *PVZ) *domain.PVZ {
	return &domain.PVZ{
		ID:               d.ID,
		RegistrationDate: d.RegistrationDate,
		City:             domain.City(d.City),
	}
}

func (PVZ) TableName() string {
	return "pvzs"
}

func (pvz PVZ) Columns() []string {
	return []string{"id", "city", "registration_date", "created_at", "updated_at", "deleted_at"}
}

func (pvz PVZ) Values() []any {
	return []any{pvz.ID, pvz.City, pvz.RegistrationDate, pvz.CreatedAt, pvz.UpdatedAt, pvz.DeletedAt}
}

func (pvz PVZ) UpdateColumns() []string {
	return []string{"id", "city", "registration_date", "updated_at"}
}

func (pvz PVZ) UpdateValues() []any {
	return []any{pvz.ID, pvz.City, pvz.RegistrationDate, time.Now()}
}
