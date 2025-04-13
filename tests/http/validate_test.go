//go:build integration.http

package http

import (
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"net/http"
)

func (s *TestSuite) TestValidationIsEnabled() {
	moderatorToken := s.createDummyLogin(api.PostDummyLoginJSONBodyRoleModerator)
	resp, err := s.client.PostPvz(s.ctx, api.PostPvzJSONRequestBody{}, createAuthRequestEditorFn(moderatorToken))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	resp, err = s.client.PostPvz(s.ctx, api.PostPvzJSONRequestBody{
		City: api.PVZCity(uuid.New().String())},
		createAuthRequestEditorFn(moderatorToken))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}
