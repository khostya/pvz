package http

import (
	"encoding/json"
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/pkg/appctx"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestProduct_PostReceptions(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  *api.PostReceptionsJSONBody
		status int
		res    any
		role   domain.Role
		mockFn func(test test, m mocks)
	}

	input := &api.PostReceptionsJSONBody{
		PvzId: uuid.New(),
	}

	tests := []test{
		{
			name:   "ok",
			input:  input,
			status: http.StatusCreated,
			res: api.Reception{
				DateTime: reception.DateTime,
				Id:       &reception.ID,
				PvzId:    reception.PvzId,
				Status:   api.ReceptionStatus(reception.Status),
			},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&reception, nil)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "forbidden",
			input:  input,
			status: http.StatusForbidden,
			res:    api.Error{Message: domain.ErrEmployeeOnly.Error()},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().Create(gomock.Any(), gomock.Any()).
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
			name:   "reception pvz error",
			input:  input,
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().Create(gomock.Any(), gomock.Any()).
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

			err := h.PostReceptions(c)
			require.NoError(t, err)

			actual := rec.Body.String()
			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.JSONEq(t, string(expected)+"\n", actual)
			require.Equal(t, tt.status, rec.Code)
		})
	}
}

func TestProduct_PvzPvzIdCloseLastReception(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  openapi_types.UUID
		status int
		res    any
		role   domain.Role
		mockFn func(test test, m mocks)
	}

	input := uuid.New()
	tests := []test{
		{
			name:   "ok",
			input:  input,
			status: http.StatusOK,
			res: api.Reception{
				DateTime: reception.DateTime,
				Id:       &reception.ID,
				PvzId:    reception.PvzId,
				Status:   api.ReceptionStatus(reception.Status),
			},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&reception, nil)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "forbidden",
			input:  input,
			status: http.StatusForbidden,
			res:    api.Error{Message: domain.ErrEmployeeOnly.Error()},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).
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
			name:   "reception pvz error",
			input:  input,
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errOops)
			},
			role: domain.UserRoleModerator,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, h, rec, mocks := newEchoCtx(t, nil)
			tt.mockFn(tt, mocks)

			if tt.role != "" {
				appctx.EchoSetRole(c, tt.role)
			}

			err := h.PostPvzPvzIdCloseLastReception(c, tt.input)
			require.NoError(t, err)

			actual := rec.Body.String()
			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.JSONEq(t, string(expected)+"\n", actual)
			require.Equal(t, tt.status, rec.Code)
		})
	}
}
