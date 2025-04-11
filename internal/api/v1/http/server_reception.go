package http

import (
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	"github.com/khostya/pvz/pkg/appctx"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s Server) PostPvzPvzIdCloseLastReception(eCtx echo.Context, pvzId openapi_types.UUID) error {
	ctx := eCtx.Request().Context()

	role, ok := appctx.EchoGetRole(eCtx)
	if !ok {
		return WriteError(eCtx, http.StatusForbidden, "role is not set")
	}
	reception, err := s.reception.CloseLastReception(ctx, dto.CloseLastReceptionParam{
		PvzID:      pvzId,
		CloserRole: role,
	})
	if isForbiddenErr(err) {
		return WriteError(eCtx, http.StatusForbidden, err.Error())
	}
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusOK, toHttpReception(reception))
}

func (s Server) PostReceptions(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var body api.PostReceptionsJSONBody
	err := eCtx.Bind(&body)
	if err != nil {
		return WriteError(eCtx, http.StatusBadRequest, err.Error())
	}

	role, ok := appctx.EchoGetRole(eCtx)
	if !ok {
		return WriteError(eCtx, http.StatusForbidden, ErrRoleIsNotSet.Error())
	}

	reception, err := s.reception.Create(ctx, dto.CreateReceptionParam{
		PvzID:       body.PvzId,
		CreatorRole: role,
	})
	if isForbiddenErr(err) {
		return WriteError(eCtx, http.StatusForbidden, err.Error())
	}
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusCreated, toHttpReception(reception))
}

func toHttpReception(r *domain.Reception) api.Reception {
	return api.Reception{
		DateTime: r.DateTime,
		PvzId:    r.PvzId,
		Id:       &r.ID,
		Status:   api.ReceptionStatus(r.Status),
	}
}
