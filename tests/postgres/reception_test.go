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

func (s *ReceptionTestSuite) TestCreateDuplicate() {
	reception := s.createReception()

	_, err := s.receptionRepo.Create(s.ctx, reception)
	s.Require().Error(repoerr.ErrDuplicate, err)
}

func (s *ReceptionTestSuite) createReception() *domain.Reception {
	pvz := NewPVZ()
	_, err := s.pvzRepo.Create(s.ctx, pvz)

	reception := NewReception(pvz.ID)
	reception, err = s.receptionRepo.Create(s.ctx, reception)
	s.Require().NoError(err)

	return reception
}

func (s *ReceptionTestSuite) TestGetByID() {
	reception := s.createReception()

	actual, err := s.receptionRepo.GetByID(s.ctx, reception.ID)
	s.Require().NoError(err)
	s.Require().EqualExportedValues(reception, actual)
}

func (s *ReceptionTestSuite) TestGetByIDNotFound() {
	_, err := s.receptionRepo.GetByID(s.ctx, uuid.New())
	s.Require().Equal(repoerr.ErrNotFound, err)
}

func (s *ReceptionTestSuite) TestUpdateLastReceptionStatus() {
	reception := s.createReception()

	s.updateLastReceptionStatusInProgress(reception, domain.ReceptionStatusClose)
	s.updateLastReceptionStatusInProgress(reception, domain.ReceptionStatusInProgress)
}

func (s *ReceptionTestSuite) TestUpdateLastReceptionStatusNotFound() {
	_, err := s.receptionRepo.UpdateReceptionStatusByID(s.ctx, uuid.New(), domain.ReceptionStatusInProgress)
	s.Require().Equal(repoerr.ErrNotFound, err)
}

func (s *ReceptionTestSuite) updateLastReceptionStatusInProgress(reception *domain.Reception, status domain.ReceptionStatus) {
	_, err := s.receptionRepo.UpdateReceptionStatusByID(s.ctx, reception.ID, status)
	s.Require().NoError(err)

	reception, err = s.receptionRepo.GetByID(s.ctx, reception.ID)
	s.Require().NoError(err)
	s.Require().EqualExportedValues(status, reception.Status)

}

func (s *ReceptionTestSuite) TestGetFirstByStatusAndPVZId() {
	_ = s.createReception()
	_ = s.createReception()
	_ = s.createReception()

	reception := s.createReception()
	actual, err := s.receptionRepo.GetFirstByStatusAndPVZId(s.ctx, reception.Status, reception.PvzId)
	s.Require().NoError(err)
	s.Require().EqualExportedValues(reception, actual)

}
