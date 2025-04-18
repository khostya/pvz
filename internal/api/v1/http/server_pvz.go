package http

import (
	"encoding/json"
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	"github.com/khostya/pvz/internal/metrics"
	"github.com/khostya/pvz/pkg/appctx"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type (
	receptions struct {
		Products  []api.Product `json:"products,omitempty"`
		Reception api.Reception `json:"reception,omitempty"`
	}

	GetPvzResponse struct {
		Pvz        api.PVZ      `json:"pvz,omitempty"`
		Receptions []receptions `json:"receptions,omitempty"`
	}
)

func (s Server) GetPvz(eCtx echo.Context, params api.GetPvzParams) error {
	keyBytes, err := json.Marshal(params)
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	key := string(keyBytes)
	if v, ok := s.getPvzResponseCache.Get(key); ok {
		return eCtx.JSON(http.StatusOK, v)
	}

	ctx := eCtx.Request().Context()

	pvzList, err := s.pvz.GetPvz(ctx, dto.GetPvzParam{
		Limit:     params.Limit,
		Page:      params.Page,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	res := make([]GetPvzResponse, 0)

	for _, pvz := range pvzList {
		resp := GetPvzResponse{
			Pvz: toHttpPVZ(pvz),
		}

		for _, r := range pvz.Receptions {
			resp.Receptions = append(resp.Receptions, receptions{
				Products:  toHttpProducts(r.Products),
				Reception: toHttpReception(r),
			})
		}

		res = append(res, resp)
	}

	s.getPvzResponseCache.Put(key, res)
	return eCtx.JSON(http.StatusOK, res)
}

func (s Server) PostPvz(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var body api.PostPvzJSONRequestBody
	err := eCtx.Bind(&body)
	if err != nil {
		return WriteError(eCtx, http.StatusBadRequest, err.Error())
	}

	role, ok := appctx.EchoGetRole(eCtx)
	if !ok {
		return WriteError(eCtx, http.StatusForbidden, ErrRoleIsNotSet.Error())
	}

	id := uuid.New()
	if body.Id != nil {
		id = *body.Id
	}

	registrationDate := time.Now()
	if body.RegistrationDate != nil {
		registrationDate = *body.RegistrationDate
	}

	pvz, err := s.pvz.Create(ctx, dto.CreatePvzParam{
		ID:               id,
		CreatorRole:      role,
		RegistrationDate: registrationDate,
		City:             domain.City(body.City),
	})
	if isForbiddenErr(err) {
		return WriteError(eCtx, http.StatusForbidden, err.Error())
	}
	if err != nil {
		return WriteError(eCtx, http.StatusInternalServerError, err.Error())
	}

	metrics.IncCreatedPVZ()
	return eCtx.JSON(http.StatusCreated, toHttpPVZ(pvz))
}

func toHttpPVZ(r *domain.PVZ) api.PVZ {
	return api.PVZ{
		City:             api.PVZCity(r.City),
		Id:               &r.ID,
		RegistrationDate: &r.RegistrationDate,
	}
}
