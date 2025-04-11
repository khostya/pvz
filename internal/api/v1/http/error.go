package http

import (
	"errors"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/labstack/echo/v4"
)

var (
	ErrRoleIsNotSet = errors.New("role is not set")
)

func WriteError(ctx echo.Context, status int, msg string) error {
	data := api.Error{Message: msg}
	return ctx.JSON(status, data)
}
