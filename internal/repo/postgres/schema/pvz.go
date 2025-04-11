package schema

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"maps"
	"time"
)

type (
	PVZ struct {
		ID               uuid.UUID `db:"pvzs.id"`
		City             string    `db:"pvzs.city"`
		RegistrationDate time.Time `db:"pvzs.registration_date"`
	}

	PvzReceptionProduct struct {
		PVZ
		Reception
		*NullableProduct
	}
)

func (p PvzReceptionProduct) Columns() []string {
	res := PVZ{}.Columns()
	res = append(res, Product{}.Columns()...)
	res = append(res, Reception{}.Columns()...)
	return res
}

func NewPVZ(d *domain.PVZ) *PVZ {
	return &PVZ{
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

func NewDomainPVZList(d []*PVZ) []*domain.PVZ {
	var res []*domain.PVZ
	for _, record := range d {
		res = append(res, NewDomainPVZ(record))
	}
	return res
}

func (PVZ) TableName() string {
	return "pvzs"
}

func (pvz PVZ) InsertColumns() []string {
	return []string{"id", "city", "registration_date"}
}

func (pvz PVZ) Columns() []string {
	return []string{"pvzs.id as \"pvzs.id\"", "pvzs.city as \"pvzs.city\"", "pvzs.registration_date as \"pvzs.registration_date\""}
}

func (pvz PVZ) Values() []any {
	return []any{pvz.ID, pvz.City, pvz.RegistrationDate}
}

func NewDomainPvzFromPvzReceptionProduct(pvzReceptionProduct []PvzReceptionProduct) []*domain.PVZ {
	var pvzs = make(map[uuid.UUID]*domain.PVZ)

	var reception = make(map[uuid.UUID]*domain.Reception)
	for _, r := range pvzReceptionProduct {
		if _, ok := pvzs[r.PVZ.ID]; !ok {
			pvzs[r.PVZ.ID] = NewDomainPVZ(&r.PVZ)
		}

		if _, ok := reception[r.Reception.ID]; !ok {
			reception[r.Reception.ID] = NewDomainReception(&r.Reception)
		}

		pvzs[r.PVZ.ID].Receptions = append(pvzs[r.PVZ.ID].Receptions, reception[r.Reception.ID])

		product := NewDomainFromNullableProduct(r.NullableProduct)
		if product != nil {
			reception[r.Reception.ID].Products = append(reception[r.Reception.ID].Products, product)
		}
	}

	var res []*domain.PVZ
	for v := range maps.Values(pvzs) {
		res = append(res, v)
	}

	return res
}

func NewDomainFromNullableProduct(n *NullableProduct) *domain.Product {
	if n == nil || n.ID == nil {
		return nil
	}

	return &domain.Product{
		ID:          *n.ID,
		DateTime:    *n.DateTime,
		Type:        domain.ProductType(*n.Type),
		ReceptionID: *n.ReceptionID,
	}
}
