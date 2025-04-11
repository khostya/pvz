package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/pkg/appctx"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
	"net/http"
	"strings"
)

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
)

type manager interface {
	ParseToken(ctx context.Context, token string) (context.Context, error)
}
type Authenticator struct {
	manager manager
}

func NewAuthenticator(manager manager) *Authenticator {
	return &Authenticator{manager: manager}
}

func (a *Authenticator) authenticatorFunc() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return authenticate(a.manager, ctx, input)
	}
}

func CreateValidatorMiddleware(authenticator *Authenticator) (echo.MiddlewareFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: authenticator.authenticatorFunc(),
			},
		})

	return validator, nil
}

func getJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get(echo.HeaderAuthorization)
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func authenticate(auth manager, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "bearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	jws, err := getJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	ctx, err = auth.ParseToken(ctx, jws)
	if err != nil {
		return fmt.Errorf("validating JWS: %w", err)
	}

	eCtx := middleware.GetEchoContext(ctx)

	appctx.SetEcho(ctx, eCtx)
	return nil
}
