//go:build integration.http

package http

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/tests/http/httpclient"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	ctx context.Context

	client api.ClientInterface
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	s.client = httpclient.NewClient()
	s.ctx = context.Background()
}

func (s *TestSuite) TestReception() {
	employeeToken := s.createDummyLogin(api.PostDummyLoginJSONBodyRoleEmployee)

	pvz, _ := s.createReception()
	_ = s.createProducts(employeeToken, 50, *pvz.Id)

	s.closeReception(employeeToken, *pvz.Id)
}

func (s *TestSuite) TestDeleteAllProducts() {
	employeeToken := s.createDummyLogin(api.PostDummyLoginJSONBodyRoleEmployee)

	pvz, _ := s.createReception()

	products := s.createProducts(employeeToken, 50, *pvz.Id)

	s.deleteProducts(employeeToken, products, *pvz.Id)
}

func (s *TestSuite) createPVZ(token string) api.PostPvzResponse {
	id := uuid.New()
	registrationDate := time.Now()

	resp, err := s.client.PostPvz(s.ctx, api.PostPvzJSONRequestBody{
		City:             api.СанктПетербург,
		Id:               &id,
		RegistrationDate: &registrationDate,
	}, createAuthRequestEditorFn(token))
	s.Require().NoError(err)

	response, err := api.ParsePostPvzResponse(resp)
	s.Require().NoError(err)

	return *response
}

func (s *TestSuite) createProducts(token string, count int, pvzID uuid.UUID) []api.Product {
	var products []api.Product
	for range count {
		products = append(products, s.createProduct(token, pvzID))
	}

	return products
}

func (s *TestSuite) deleteProducts(token string, products []api.Product, pvzID uuid.UUID) {
	for range products {
		lastProduct, err := s.client.PostPvzPvzIdDeleteLastProduct(s.ctx, pvzID, createAuthRequestEditorFn(token))
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.StatusOK, lastProduct.StatusCode)
	}
}

func (s *TestSuite) createProduct(token string, pvzID uuid.UUID) api.Product {
	body := api.PostProductsJSONRequestBody{
		PvzId: pvzID,
		Type:  api.PostProductsJSONBodyTypeОбувь,
	}

	resp, err := s.client.PostProducts(s.ctx, body, createAuthRequestEditorFn(token))
	s.Require().NoError(err)

	response, err := api.ParsePostProductsResponse(resp)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusCreated, response.StatusCode(), string(response.Body))
	s.Require().NotNil(response.JSON201)

	return *response.JSON201
}

func (s *TestSuite) closeReception(token string, pvzID uuid.UUID) api.Reception {
	resp, err := s.client.PostPvzPvzIdCloseLastReception(s.ctx, pvzID, createAuthRequestEditorFn(token))
	s.Require().NoError(err)

	response, err := api.ParsePostPvzPvzIdCloseLastReceptionResponse(resp)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, response.StatusCode(), string(response.Body))
	s.Require().NotNil(response.JSON200)

	return *response.JSON200
}

func (s *TestSuite) createReception() (api.PVZ, api.Reception) {
	moderatorToken := s.createDummyLogin(api.PostDummyLoginJSONBodyRoleModerator)
	employeeToken := s.createDummyLogin(api.PostDummyLoginJSONBodyRoleEmployee)

	createPvzResp := s.createPVZ(moderatorToken)
	s.Require().Equal(http.StatusCreated, createPvzResp.StatusCode())
	s.Require().NotNil(createPvzResp.JSON201)

	pvz := *createPvzResp.JSON201

	resp, err := s.client.PostReceptions(s.ctx, api.PostReceptionsJSONRequestBody{
		PvzId: *pvz.Id,
	}, createAuthRequestEditorFn(employeeToken))
	s.Require().NoError(err)

	response, err := api.ParsePostReceptionsResponse(resp)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusCreated, response.StatusCode(), string(response.Body))
	s.Require().NotNil(response.JSON201)

	return pvz, *response.JSON201
}

func createAuthRequestEditorFn(token string) api.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
		return nil
	}
}
