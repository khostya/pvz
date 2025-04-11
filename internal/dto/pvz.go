package dto

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

const (
	defaultLimit = 10
)

type (
	CreatePvzParam struct {
		ID               uuid.UUID   `json:"id" validate:"required"`
		City             domain.City `json:"city" validate:"required"`
		RegistrationDate time.Time   `json:"registrationDate" validate:"required"`
		CreatorRole      domain.Role `json:"creatorRole" validate:"required"`
	}

	GetPvzParam struct {
		StartDate *time.Time `json:"startDate"`
		EndDate   *time.Time `json:"endDate"`
		Page      *int       `json:"page"`
		Limit     *int       `json:"limit"`
	}
)

func (p GetPvzParam) Offset() uint64 {
	if p.Page == nil {
		return uint64(p.Count())
	}
	return uint64(*p.Page-1) * p.Count()
}

func (p GetPvzParam) Count() uint64 {
	if p.Limit == nil {
		return defaultLimit
	}
	return uint64(*p.Limit)
}
