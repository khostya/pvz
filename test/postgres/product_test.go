//go:build integration

package postgres

import (
	"context"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/repo/postgres"
	"github.com/khostya/pvz/pkg/postgres/repoerr"
	"github.com/khostya/pvz/pkg/postgres/transactor"
	"github.com/stretchr/testify/require"
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

	reception := NewReception(pvz.ID)
	reception, err = s.receptionRepo.Create(s.ctx, reception)
	require.NoError(s.T(), err)

	product := NewProduct(reception.ID, pvz.ID)
	product, err = s.productRepo.Create(s.ctx, product)
	require.NoError(s.T(), err)
	return product
}

func (s *ProductTestSuite) TestGetByID() {
	product := s.createProduct()

	actual, err := s.productRepo.GetByID(s.ctx, product.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), product, actual)
}

func (s *ProductTestSuite) TestDeleteByDateTIme() {
	product := s.createProduct()

	err := s.productRepo.DeleteLastByDateTime(s.ctx, product.ReceptionID)
	require.NoError(s.T(), err)

	_, err = s.productRepo.GetByID(s.ctx, product.ID)
	require.Equal(s.T(), repoerr.ErrNotFound, err)
}

func (s *ProductTestSuite) TestDeleteLastByDateTIme() {
	product := s.createProduct()

	product2 := NewProduct(product.ReceptionID, product.PvzID)
	product2, err := s.productRepo.Create(s.ctx, product2)
	require.NoError(s.T(), err)

	err = s.productRepo.DeleteLastByDateTime(s.ctx, product.ReceptionID)
	require.NoError(s.T(), err)

	actual, err := s.productRepo.GetByID(s.ctx, product.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), product, actual)

	_, err = s.productRepo.GetByID(s.ctx, product2.ID)
	require.Equal(s.T(), repoerr.ErrNotFound, err)
}
