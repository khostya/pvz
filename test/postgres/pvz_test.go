//go:build integration

package postgres

import (
	"context"
	"github.com/khostya/pvz/internal/repo/postgres"
	"github.com/khostya/pvz/pkg/postgres/transactor"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PVZsTestSuite struct {
	suite.Suite
	ctx        context.Context
	pvzRepo    *postgres.PvzRepo
	transactor *transactor.TransactionManager
}

func TestPVZs(t *testing.T) {
	suite.Run(t, new(PVZsTestSuite))
}

func (s *PVZsTestSuite) SetupSuite() {
	s.transactor = db.GetTransactionManager()
	s.pvzRepo = postgres.NewPvzRepo(s.transactor)
	s.ctx = context.Background()
}

func (s *PVZsTestSuite) TestCreate() {
	pvz := NewPVZ()

	_, err := s.pvzRepo.Create(s.ctx, pvz)
	require.NoError(s.T(), err)
}

func (s *PVZsTestSuite) TestGetByID() {
	pvz := NewPVZ()

	pvz, err := s.pvzRepo.Create(s.ctx, pvz)
	require.NoError(s.T(), err)

	actual, err := s.pvzRepo.GetByID(s.ctx, pvz.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), pvz, actual)
}
