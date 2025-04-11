//go:build integration

package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	"github.com/khostya/pvz/internal/repo/postgres"
	"github.com/khostya/pvz/pkg/postgres/repoerr"
	"github.com/khostya/pvz/pkg/postgres/transactor"
	"github.com/stretchr/testify/suite"
	"sort"
	"testing"
	"time"
)

var (
	page  = 1
	limit = 100
)

type PVZsTestSuite struct {
	suite.Suite
	ctx           context.Context
	pvzRepo       *postgres.PvzRepo
	transactor    *transactor.TransactionManager
	productRepo   *postgres.ProductRepo
	receptionRepo *postgres.ReceptionRepo
}

func TestPVZs(t *testing.T) {
	suite.Run(t, new(PVZsTestSuite))
}

func (s *PVZsTestSuite) SetupSuite() {
	s.transactor = db.GetTransactionManager()
	s.pvzRepo = postgres.NewPvzRepo(s.transactor)
	s.receptionRepo = postgres.NewReceptionRepo(s.transactor)
	s.productRepo = postgres.NewProductRepo(s.transactor)
	s.ctx = context.Background()
}

func (s *PVZsTestSuite) TestCreate() {
	_ = s.create()
}

func (s *PVZsTestSuite) TestCreateDuplicate() {
	pvz := s.create()

	_, err := s.pvzRepo.Create(s.ctx, pvz)
	s.Require().Equal(repoerr.ErrDuplicate, err)
}

func (s *PVZsTestSuite) create() *domain.PVZ {
	pvz := NewPVZ()

	pvz, err := s.pvzRepo.Create(s.ctx, pvz)
	s.Require().NoError(err)

	return pvz
}

func (s *PVZsTestSuite) TestGetByID() {
	pvz := s.create()

	actual, err := s.pvzRepo.GetByID(s.ctx, pvz.ID)
	s.Require().NoError(err)
	s.Require().EqualExportedValues(pvz, actual)
}

func (s *PVZsTestSuite) TestGetByIDNotFound() {
	_, err := s.pvzRepo.GetByID(s.ctx, uuid.New())
	s.Require().Equal(repoerr.ErrNotFound, err)
}

func (s *PVZsTestSuite) TestGetAll() {
	truncate()

	pvz := s.create()

	actual, err := s.pvzRepo.GetAllPVZList(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(actual, 1)
	s.Require().EqualExportedValues(pvz, actual[0])
}

func (s *PVZsTestSuite) TestGetPvz() {
	truncate()

	pvz, reception := s.createFullPvz()
	s.checkGetPvz([]*domain.PVZ{pvz}, reception.DateTime, reception.DateTime)
}

func (s *PVZsTestSuite) TestGetPvzManyPvz() {
	truncate()

	var res []*domain.PVZ

	pvz, reception := s.createFullPvz()
	res = append(res, pvz)
	startDate := reception.DateTime

	for range 10 {
		pvz, _ = s.createFullPvz()
		res = append(res, pvz)

		pvz, _ = s.createFullPvzWithoutProducts()
		res = append(res, pvz)
	}

	pvz, reception = s.createFullPvzWithoutProducts()
	res = append(res, pvz)

	s.checkGetPvz(res, startDate, reception.DateTime)
}

func (s *PVZsTestSuite) TestGetPvzWithoutProducts() {
	truncate()

	pvz, reception := s.createFullPvzWithoutProducts()

	s.checkGetPvz([]*domain.PVZ{pvz}, reception.DateTime, reception.DateTime)
}

func (s *PVZsTestSuite) checkGetPvz(expected []*domain.PVZ, startDate, endDate time.Time) {
	actual, err := s.pvzRepo.GetPVZ(s.ctx, dto.GetPvzParam{
		Page:      &page,
		Limit:     &limit,
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	s.Require().NoError(err)
	s.Require().Len(actual, len(expected))

	sort.Slice(expected, func(i, j int) bool {
		return expected[i].ID.String() < expected[j].ID.String()
	})
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].ID.String() < actual[j].ID.String()
	})

	s.Require().EqualExportedValues(expected, actual)
}

func (s *PVZsTestSuite) createFullPvz() (*domain.PVZ, *domain.Reception) {
	pvz, reception := s.createFullPvzWithoutProducts()

	product := NewProduct(reception.ID)
	product, err := s.productRepo.Create(s.ctx, product)
	s.Require().NoError(err)

	reception.Products = append(reception.Products, product)

	return pvz, reception
}

func (s *PVZsTestSuite) createFullPvzWithoutProducts() (*domain.PVZ, *domain.Reception) {
	pvz := NewPVZ()
	_, err := s.pvzRepo.Create(s.ctx, pvz)
	s.Require().NoError(err)

	reception := NewReception(pvz.ID)
	reception, err = s.receptionRepo.Create(s.ctx, reception)
	s.Require().NoError(err)
	pvz.Receptions = append(pvz.Receptions, reception)

	return pvz, reception
}
