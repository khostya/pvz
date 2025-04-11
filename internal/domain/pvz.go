package domain

import (
	"github.com/google/uuid"
	"time"
)

const (
	CityMoscow = "Москва"
)

type (
	City string

	PVZ struct {
		ID               uuid.UUID `json:"id"`
		RegistrationDate time.Time `json:"registrationDate"`
		City             City      `json:"city"`

		Receptions []*Reception `json:"receptions"`
	}
)
