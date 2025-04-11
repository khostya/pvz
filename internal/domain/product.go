package domain

import (
	"github.com/google/uuid"
	"time"
)

const (
	ProductTypeShoes ProductType = "обувь"
)

type (
	ProductType string

	Product struct {
		ID          uuid.UUID   `json:"id"`
		DateTime    time.Time   `json:"dateTime"`
		Type        ProductType `json:"type"`
		PvzID       uuid.UUID   `json:"pvzId"`
		ReceptionID uuid.UUID   `json:"receptionId"`
	}
)
