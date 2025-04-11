//go:build integration

package postgres

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

func NewUser() *domain.User {
	return &domain.User{
		ID:       uuid.New(),
		Email:    gofakeit.Email(),
		Password: uuid.New().String(),
		Role:     domain.UserRoleModerator,
	}
}

func NewPVZ() *domain.PVZ {
	return &domain.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             domain.CityMoscow,
	}
}

func NewReception(pvzID uuid.UUID) *domain.Reception {
	return &domain.Reception{
		ID:       uuid.New(),
		PvzId:    pvzID,
		DateTime: time.Now(),
		Status:   domain.ReceptionStatusClose,
	}
}

func NewProduct(receptionID, pvzID uuid.UUID) *domain.Product {
	return &domain.Product{
		ID:          uuid.New(),
		ReceptionID: receptionID,
		PvzID:       pvzID,
		DateTime:    time.Now(),
		Type:        domain.ProductTypeShoes,
	}
}
