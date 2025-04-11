package dto

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
)

type (
	CreateReceptionParam struct {
		CreatorRole domain.Role `json:"creatorRole" validate:"required"`
		PvzID       uuid.UUID   `json:"pvzId" validate:"required"`
	}

	CloseLastReceptionParam struct {
		PvzID      uuid.UUID   `json:"pvzId" validate:"required"`
		CloserRole domain.Role `json:"closerRole" validate:"required"`
	}

	DeleteLastReceptionParam struct {
		PvzID       uuid.UUID   `json:"pvzId" validate:"required"`
		DeleterRole domain.Role `json:"deleterRole" validate:"required"`
	}
)
