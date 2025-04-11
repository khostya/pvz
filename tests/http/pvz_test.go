//go:build integration.http

package http

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/tests/http/httpclient"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type PvzTestSuite struct {
	suite.Suite
	ctx context.Context

	client api.ClientInterface
}

func TestPVZ(t *testing.T) {
	suite.Run(t, new(PvzTestSuite))
}

func (s *PvzTestSuite) SetupSuite() {
	s.client = httpclient.NewClient()
	s.ctx = context.Background()
}

func (s *PvzTestSuite) Test() {
	moderatorToken := s.createDummyLogin(api.PostDummyLoginJSONBodyRoleModerator)
	employeeToken := s.createDummyLogin(api.PostDummyLoginJSONBodyRoleEmployee)

	createPvzResp := s.createPVZ(moderatorToken)
	s.Require().Equal(http.StatusCreated, createPvzResp.StatusCode())
	s.Require().NotNil(createPvzResp.JSON201)

	pvz := *createPvzResp.JSON201

	_ = s.createReception(*pvz.Id, employeeToken)
	_ = s.createProducts(employeeToken, 50, *pvz.Id)

	s.closeReception(employeeToken, *pvz.Id)
}

func (s *PvzTestSuite) TestAuthorizationIsEnabled() {
	createPvzResp := s.createPVZ("fsd")
	s.Require().Equal(http.StatusForbidden, createPvzResp.StatusCode())
}

func (s *PvzTestSuite) TestValidationIsEnabled() {
	moderatorToken := s.createDummyLogin(api.PostDummyLoginJSONBodyRoleModerator)
	resp, err := s.client.PostPvz(s.ctx, api.PostPvzJSONRequestBody{}, CreateAuthRequestEditorFn(moderatorToken))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	resp, err = s.client.PostPvz(s.ctx, api.PostPvzJSONRequestBody{
		City: api.PVZCity(uuid.New().String())},
		CreateAuthRequestEditorFn(moderatorToken))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *PvzTestSuite) createPVZ(token string) api.PostPvzResponse {
	id := uuid.New()
	registrationDate := time.Now()

	resp, err := s.client.PostPvz(s.ctx, api.PostPvzJSONRequestBody{
		City:             api.СанктПетербург,
		Id:               &id,
		RegistrationDate: &registrationDate,
	}, CreateAuthRequestEditorFn(token))
	s.Require().NoError(err)

	response, err := api.ParsePostPvzResponse(resp)
	s.Require().NoError(err)

	return *response
}

func (s *PvzTestSuite) createProducts(token string, count int, pvzID uuid.UUID) []api.Product {
	body := api.PostProductsJSONRequestBody{
		PvzId: pvzID,
		Type:  api.PostProductsJSONBodyTypeОбувь,
	}

	var products []api.Product
	for range count {
		resp, err := s.client.PostProducts(s.ctx, body, CreateAuthRequestEditorFn(token))
		s.Require().NoError(err)

		response, err := api.ParsePostProductsResponse(resp)
		s.Require().NoError(err)

		s.Require().Equal(http.StatusCreated, response.StatusCode(), string(response.Body))
		s.Require().NotNil(response.JSON201)

		products = append(products, *response.JSON201)
	}

	return products
}

func (s *PvzTestSuite) closeReception(token string, pvzID uuid.UUID) api.Reception {
	resp, err := s.client.PostPvzPvzIdCloseLastReception(s.ctx, pvzID, CreateAuthRequestEditorFn(token))
	s.Require().NoError(err)

	response, err := api.ParsePostPvzPvzIdCloseLastReceptionResponse(resp)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, response.StatusCode(), string(response.Body))
	s.Require().NotNil(response.JSON200)

	return *response.JSON200
}

func (s *PvzTestSuite) createReception(pvzID openapi_types.UUID, token string) api.Reception {
	resp, err := s.client.PostReceptions(s.ctx, api.PostReceptionsJSONRequestBody{
		PvzId: pvzID,
	}, CreateAuthRequestEditorFn(token))
	s.Require().NoError(err)

	response, err := api.ParsePostReceptionsResponse(resp)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusCreated, response.StatusCode(), string(response.Body))
	s.Require().NotNil(response.JSON201)

	return *response.JSON201
}

func (s *PvzTestSuite) createDummyLogin(role api.PostDummyLoginJSONBodyRole) api.Token {
	resp, err := s.client.PostDummyLogin(s.ctx, api.PostDummyLoginJSONRequestBody{
		Role: role,
	})
	s.Require().NoError(err)

	response, err := api.ParsePostDummyLoginResponse(resp)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, response.StatusCode(), string(response.Body))
	s.Require().NotNil(response.JSON200)

	return *response.JSON200
}

func CreateAuthRequestEditorFn(token string) api.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
		return nil
	}
}
