//go:build integration.http

package http

import (
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"net/http"
)

func (s *TestSuite) TestAuthorizationIsEnabled() {
	createPvzResp := s.createPVZ("fsd")
	s.Require().Equal(http.StatusForbidden, createPvzResp.StatusCode())
}

func (s *TestSuite) createDummyLogin(role api.PostDummyLoginJSONBodyRole) api.Token {
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
