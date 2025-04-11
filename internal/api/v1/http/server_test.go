package http

import (
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

var (
	pvz = domain.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             domain.CityMoscow,
		Receptions:       []*domain.Reception{&reception},
	}

	product = domain.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        domain.ProductTypeShoes,
		PvzID:       uuid.New(),
		ReceptionID: uuid.New(),
	}

	reception = domain.Reception{
		ID:       uuid.New(),
		PvzId:    uuid.New(),
		DateTime: time.Now(),
		Status:   domain.ReceptionStatus(api.Close),
		Products: []*domain.Product{&product},
	}

	user = &domain.User{
		Email:    "khostya.konsantin@gmail.com",
		Password: "bla-bla-bla",
		ID:       uuid.New(),
		Role:     domain.UserRoleModerator,
	}
)
