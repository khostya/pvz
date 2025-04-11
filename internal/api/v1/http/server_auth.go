package http

import (
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s Server) PostDummyLogin(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var body api.PostDummyLoginJSONBody
	err := eCtx.Bind(&body)
	if err != nil {
		return WriteError(eCtx, http.StatusBadRequest, err.Error())
	}

	token, err := s.auth.DummyLogin(ctx, dto.DummyLoginUserParam{
		Role: domain.Role(body.Role),
	})
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusOK, token)
}

func (s Server) PostLogin(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var body api.PostLoginJSONBody
	err := eCtx.Bind(&body)
	if err != nil {
		return WriteError(eCtx, http.StatusBadRequest, err.Error())
	}

	token, err := s.auth.Login(ctx, dto.LoginUserParam{
		Email:    string(body.Email),
		Password: body.Password,
	})
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusOK, token)
}

func (s Server) PostRegister(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var body api.PostRegisterJSONBody
	err := eCtx.Bind(&body)
	if err != nil {
		return WriteError(eCtx, http.StatusBadRequest, err.Error())
	}

	user, err := s.auth.Register(ctx, dto.RegisterUserParam{
		Email:    string(body.Email),
		Password: body.Password,
		Role:     domain.Role(body.Role),
	})
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusCreated, api.User{
		Email: openapi_types.Email(user.Email),
		Id:    &user.ID,
		Role:  api.UserRole(user.Role),
	})
}
