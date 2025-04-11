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

type UsersTestSuite struct {
	suite.Suite
	ctx        context.Context
	userRepo   *postgres.UserRepo
	transactor *transactor.TransactionManager
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

func (s *UsersTestSuite) SetupSuite() {
	s.transactor = db.GetTransactionManager()
	s.userRepo = postgres.NewUserRepo(s.transactor)
	s.ctx = context.Background()
}

func (s *UsersTestSuite) TestCreate() {
	_ = s.create()
}

func (s *UsersTestSuite) TestCreateDuplicate() {
	user := s.create()

	_, err := s.userRepo.Create(s.ctx, user)
	s.Require().Equal(repoerr.ErrDuplicate, err)
}

func (s *UsersTestSuite) TestGetByID() {
	user := NewUser()

	user, err := s.userRepo.Create(s.ctx, user)
	s.Require().NoError(err)

	actual, err := s.userRepo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Require().EqualExportedValues(user, actual)
}

func (s *UsersTestSuite) TestGetByEmailID() {
	user := s.create()

	actual, err := s.userRepo.GetByEmail(s.ctx, user.Email)
	s.Require().NoError(err)
	s.Require().EqualExportedValues(user, actual)
}

func (s *UsersTestSuite) TestGetByIDNotFound() {
	_, err := s.userRepo.GetByID(s.ctx, uuid.New())
	s.Require().Equal(repoerr.ErrNotFound, err)
}

func (s *UsersTestSuite) TestGetByEmailNotFound() {
	_, err := s.userRepo.GetByEmail(s.ctx, uuid.New().String())
	s.Require().Equal(repoerr.ErrNotFound, err)
}

func (s *UsersTestSuite) create() *domain.User {
	user := NewUser()

	user, err := s.userRepo.Create(s.ctx, user)
	s.Require().NoError(err)

	return user
}
