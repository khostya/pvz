package http

import (
	"errors"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	"github.com/khostya/pvz/internal/metrics"
	"github.com/khostya/pvz/pkg/appctx"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s Server) PostProducts(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	role, ok := appctx.EchoGetRole(eCtx)
	if !ok {
		return WriteError(eCtx, http.StatusForbidden, ErrRoleIsNotSet.Error())
	}

	var body api.PostProductsJSONBody
	err := eCtx.Bind(&body)
	if err != nil {
		return WriteError(eCtx, http.StatusBadRequest, err.Error())
	}

	product, err := s.product.Create(ctx, dto.CreateProductParam{
		CreatorRole: role,
		PvzID:       body.PvzId,
		Type:        domain.ProductType(body.Type),
	})
	if isForbiddenErr(err) {
		return WriteError(eCtx, http.StatusForbidden, err.Error())
	}
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	metrics.IncCreatedProducts()
	return eCtx.JSON(http.StatusCreated, toHttpProduct(product))
}

func (s Server) PostPvzPvzIdDeleteLastProduct(eCtx echo.Context, pvzId openapi_types.UUID) error {
	ctx := eCtx.Request().Context()

	role, ok := appctx.EchoGetRole(eCtx)
	if !ok {
		return WriteError(eCtx, http.StatusForbidden, ErrRoleIsNotSet.Error())
	}

	err := s.reception.DeleteLastProduct(ctx, dto.DeleteLastReceptionParam{
		PvzID:       pvzId,
		DeleterRole: role,
	})
	if isForbiddenErr(err) {
		return WriteError(eCtx, http.StatusForbidden, err.Error())
	}
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusOK, nil)
}

func isForbiddenErr(err error) bool {
	return errors.Is(err, domain.ErrEmployeeOnly) || errors.Is(err, domain.ErrModeratorOnly)
}

func toHttpProduct(p *domain.Product) api.Product {
	return api.Product{
		DateTime:    &p.DateTime,
		Id:          &p.ID,
		ReceptionId: p.ReceptionID,
		Type:        api.ProductType(p.Type),
	}
}

func toHttpProducts(p []*domain.Product) []api.Product {
	var res = make([]api.Product, len(p))
	for i := range p {
		res[i] = toHttpProduct(p[i])
	}
	return res
}
