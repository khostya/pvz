//go:build integration

package postgres

import (
	"context"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/repo/postgres"
	"github.com/khostya/pvz/pkg/postgres/transactor"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ReceptionTestSuite struct {
	suite.Suite
	ctx           context.Context
	pvzRepo       *postgres.PvzRepo
	receptionRepo *postgres.ReceptionRepo
	transactor    *transactor.TransactionManager
}

func TestReception(t *testing.T) {
	suite.Run(t, new(ReceptionTestSuite))
}

func (s *ReceptionTestSuite) SetupSuite() {
	s.transactor = db.GetTransactionManager()
	s.pvzRepo = postgres.NewPvzRepo(s.transactor)
	s.receptionRepo = postgres.NewReceptionRepo(s.transactor)

	s.ctx = context.Background()
}

func (s *ReceptionTestSuite) TestCreate() {
	_ = s.createReception()
}

func (s *ReceptionTestSuite) createReception() *domain.Reception {
	pvz := NewPVZ()
	_, err := s.pvzRepo.Create(s.ctx, pvz)

	reception := NewReception(pvz.ID)
	reception, err = s.receptionRepo.Create(s.ctx, reception)
	require.NoError(s.T(), err)

	return reception
}

func (s *ReceptionTestSuite) TestGetByID() {
	reception := s.createReception()

	actual, err := s.receptionRepo.GetByID(s.ctx, reception.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), reception, actual)
}
