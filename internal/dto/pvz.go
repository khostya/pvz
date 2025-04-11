package dto

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
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
