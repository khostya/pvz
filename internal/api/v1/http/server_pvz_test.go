package http

import (
	"encoding/json"
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/pkg/appctx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
	"time"
)

func TestProduct_PostPvz(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  *api.PostPvzJSONRequestBody
		status int
		res    any
		role   domain.Role
		mockFn func(test test, m mocks)
	}

	id := uuid.New()
	registrationDate := time.Now()
	input := &api.PostPvzJSONRequestBody{
		City:             api.СанктПетербург,
		Id:               &id,
		RegistrationDate: &registrationDate,
	}

	tests := []test{
		{
			name:   "ok",
			input:  input,
			status: http.StatusCreated,
			res: api.PVZ{
				City:             input.City,
				Id:               input.Id,
				RegistrationDate: &registrationDate,
			},
			mockFn: func(test test, m mocks) {
				m.pvz.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.PVZ{
						ID:               *test.input.Id,
						RegistrationDate: *test.input.RegistrationDate,
						City:             domain.City(test.input.City),
					}, nil)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "forbidden",
			input:  input,
			status: http.StatusForbidden,
			res:    api.Error{Message: domain.ErrEmployeeOnly.Error()},
			mockFn: func(test test, m mocks) {
				m.pvz.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, domain.ErrEmployeeOnly)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "role is not set",
			input:  input,
			status: http.StatusForbidden,
			res:    api.Error{Message: ErrRoleIsNotSet.Error()},
			mockFn: func(test test, m mocks) {

			},
		},
		{
			name:   "create pvz error",
			input:  input,
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.pvz.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errOops)
			},
			role: domain.UserRoleModerator,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, h, rec, mocks := newEchoCtx(t, tt.input)
			tt.mockFn(tt, mocks)

			if tt.role != "" {
				appctx.EchoSetRole(c, tt.role)
			}

			err := h.PostPvz(c)
			require.NoError(t, err)

			actual := rec.Body.String()
			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.JSONEq(t, string(expected)+"\n", actual)
			require.Equal(t, tt.status, rec.Code)
		})
	}
}

func TestProduct_GetPvz(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  api.GetPvzParams
		status int
		res    any
		mockFn func(test test, m mocks)
	}

	input := api.GetPvzParams{}
	pvzResponse := []GetPvzResponse{
		{
			Pvz: api.PVZ{
				Id:               &pvz.ID,
				City:             api.PVZCity(pvz.City),
				RegistrationDate: &pvz.RegistrationDate,
			},
			Receptions: []receptions{
				{
					Reception: toHttpReception(&reception),
					Products:  []api.Product{toHttpProduct(&product)},
				},
			},
		},
	}
	tests := []test{
		{
			name:   "ok",
			input:  input,
			status: http.StatusOK,
			res:    pvzResponse,
			mockFn: func(test test, m mocks) {
				m.getPvzResponseCache.EXPECT().Get(gomock.Any()).
					Times(1).
					Return(nil, false)
				m.pvz.EXPECT().GetPvz(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]*domain.PVZ{&pvz}, nil)
				m.getPvzResponseCache.EXPECT().Put(gomock.Any(), test.res).
					Times(1).
					Return()
			},
		},
		{
			name:   "cache hit",
			input:  input,
			status: http.StatusOK,
			res:    pvzResponse,
			mockFn: func(test test, m mocks) {
				m.getPvzResponseCache.EXPECT().Get(gomock.Any()).
					Times(1).
					Return(test.res, true)
			},
		},
		{
			name:   "err get pvz",
			input:  input,
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.getPvzResponseCache.EXPECT().Get(gomock.Any()).
					Times(1).
					Return(nil, false)
				m.pvz.EXPECT().GetPvz(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errOops)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, h, rec, mocks := newEchoCtx(t, nil)
			tt.mockFn(tt, mocks)

			err := h.GetPvz(c, tt.input)
			require.NoError(t, err)

			actual := rec.Body.String()
			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.JSONEq(t, string(expected)+"\n", actual)
			require.Equal(t, tt.status, rec.Code)
		})
	}
}
