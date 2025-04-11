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
	user := NewUser()

	_, err := s.userRepo.Create(s.ctx, user)
	require.NoError(s.T(), err)
}

func (s *UsersTestSuite) TestGetByID() {
	user := NewUser()

	user, err := s.userRepo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	actual, err := s.userRepo.GetByID(s.ctx, user.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), user, actual)
}
