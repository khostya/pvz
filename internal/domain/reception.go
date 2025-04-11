package domain

import (
	"github.com/google/uuid"
	"time"
)

const (
	ReceptionStatusClose      ReceptionStatus = "close"
	ReceptionStatusInProgress ReceptionStatus = "in_progress"
)

type (
	ReceptionStatus string

	Reception struct {
		ID       uuid.UUID       `json:"id"`
		PvzId    uuid.UUID       `json:"pvzId"`
		DateTime time.Time       `json:"dateTime"`
		Status   ReceptionStatus `json:"status"`

		Products []*Product `json:"products"`
	}
)
