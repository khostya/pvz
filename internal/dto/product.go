package dto

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
)

type CreateProductParam struct {
	CreatorRole domain.Role        `json:"creatorRole" validate:"required"`
	PvzID       uuid.UUID          `json:"pvzId" validate:"required"`
	Type        domain.ProductType `json:"type" validate:"required"`
}
