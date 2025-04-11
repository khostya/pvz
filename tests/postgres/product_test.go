//go:build integration

package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/repo/postgres"
	"github.com/khostya/pvz/pkg/postgres/repoerr"
	"github.com/khostya/pvz/pkg/postgres/transactor"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ProductTestSuite struct {
	suite.Suite
	ctx context.Context

	pvzRepo       *postgres.PvzRepo
	receptionRepo *postgres.ReceptionRepo
	productRepo   *postgres.ProductRepo

	transactor *transactor.TransactionManager
}

func TestProduct(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (s *ProductTestSuite) SetupSuite() {
	s.transactor = db.GetTransactionManager()
	s.pvzRepo = postgres.NewPvzRepo(s.transactor)
	s.receptionRepo = postgres.NewReceptionRepo(s.transactor)
	s.productRepo = postgres.NewProductRepo(s.transactor)
	s.ctx = context.Background()
}

func (s *ProductTestSuite) TestCreate() {
	_ = s.createProduct()
}

func (s *ProductTestSuite) createProduct() *domain.Product {
	pvz := NewPVZ()
	_, err := s.pvzRepo.Create(s.ctx, pvz)
	s.Require().NoError(err)

	reception := NewReception(pvz.ID)
	reception, err = s.receptionRepo.Create(s.ctx, reception)
	s.Require().NoError(err)

	product := NewProduct(reception.ID)
	product, err = s.productRepo.Create(s.ctx, product)
	s.Require().NoError(err)

	return product
}

func (s *ProductTestSuite) TestCreateDuplicate() {
	product := s.createProduct()

	_, err := s.productRepo.Create(s.ctx, product)
	s.Require().Equal(repoerr.ErrDuplicate, err)
}

func (s *ProductTestSuite) TestGetByID() {
	product := s.createProduct()

	actual, err := s.productRepo.GetByID(s.ctx, product.ID)
	s.Require().NoError(err)
	s.Require().EqualExportedValues(product, actual)
}

func (s *ProductTestSuite) TestGetByIDNotFound() {
	_, err := s.productRepo.GetByID(s.ctx, uuid.New())
	s.Require().Equal(repoerr.ErrNotFound, err)
}

func (s *ProductTestSuite) TestDeleteByDateTIme() {
	product := s.createProduct()

	err := s.productRepo.DeleteLastByDateTime(s.ctx, product.ReceptionID)
	s.Require().NoError(err)

	_, err = s.productRepo.GetByID(s.ctx, product.ID)
	s.Require().Equal(repoerr.ErrNotFound, err)
}

func (s *ProductTestSuite) TestDeleteByDateTImeNotFound() {
	err := s.productRepo.DeleteLastByDateTime(s.ctx, uuid.New())
	s.Require().Equal(repoerr.ErrNotFound, err)
}

func (s *ProductTestSuite) TestDeleteLastByDateTIme() {
	truncate()

	product := s.createProduct()

	product2 := NewProduct(product.ReceptionID)
	product2, err := s.productRepo.Create(s.ctx, product2)
	s.Require().NoError(err)

	err = s.productRepo.DeleteLastByDateTime(s.ctx, product.ReceptionID)
	s.Require().NoError(err)

	actual, err := s.productRepo.GetByID(s.ctx, product2.ID)
	s.Require().NoError(err)
	s.Require().EqualExportedValues(product2, actual)

	_, err = s.productRepo.GetByID(s.ctx, product.ID)
	s.Require().Equal(repoerr.ErrNotFound, err)
}
