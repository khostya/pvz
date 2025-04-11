//go:build integration.grpc

package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	pvz_v1 "github.com/khostya/pvz/pkg/api/v1/proto"
	"github.com/khostya/pvz/tests/grpc/grpcclient"
	"github.com/khostya/pvz/tests/http/httpclient"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type PvzTestSuite struct {
	suite.Suite
	ctx context.Context

	httpClient api.ClientInterface
	grpcClient *grpcclient.Client
}

func TestPVZ(t *testing.T) {
	suite.Run(t, new(PvzTestSuite))
}

func (s *PvzTestSuite) SetupSuite() {
	s.httpClient = httpclient.NewClient()
	s.grpcClient = grpcclient.NewClient()
	s.ctx = context.Background()
}

func (s *PvzTestSuite) TearDownSuite() {
	err := s.grpcClient.Close()
	s.Require().NoError(err)
}

func (s *PvzTestSuite) TestGetPvz() {
	pvz := s.createPVZ()

	resp, err := s.grpcClient.GetPVZList(s.ctx, &pvz_v1.GetPVZListRequest{})
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Len(resp.Pvzs, 1)

	grpcPvz := resp.Pvzs[0]

	id := uuid.MustParse(grpcPvz.Id)
	registrationDate := grpcPvz.RegistrationDate.AsTime()

	s.Require().Equal(pvz, api.PVZ{
		City:             api.PVZCity(grpcPvz.City),
		Id:               &id,
		RegistrationDate: &registrationDate,
	})
}

func (s *PvzTestSuite) createPVZ() api.PVZ {
	resp, err := s.httpClient.PostDummyLogin(s.ctx, api.PostDummyLoginJSONRequestBody{
		Role: api.PostDummyLoginJSONBodyRoleModerator,
	})
	s.Require().NoError(err)

	response, err := api.ParsePostDummyLoginResponse(resp)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, response.StatusCode(), string(response.Body))
	s.Require().NotNil(response.JSON200)

	token := *response.JSON200

	id := uuid.New()
	registrationDate := time.Now()

	resp, err = s.httpClient.PostPvz(s.ctx, api.PostPvzJSONRequestBody{
		City:             api.СанктПетербург,
		Id:               &id,
		RegistrationDate: &registrationDate,
	}, s.createAuthRequestEditorFn(token))
	s.Require().NoError(err)

	pvzResponse, err := api.ParsePostPvzResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(pvzResponse.JSON201)

	return *pvzResponse.JSON201
}

func (s *PvzTestSuite) createAuthRequestEditorFn(token string) api.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
		return nil
	}
}
