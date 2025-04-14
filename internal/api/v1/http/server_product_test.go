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

func TestProduct_PostProducts(t *testing.T) {
	t.Parallel()

	type test struct {
		name  string
		input *api.PostProductsJSONBody

		status int
		res    any
		role   domain.Role

		mockFn func(test test, m mocks)
	}

	input := &api.PostProductsJSONBody{
		PvzId: uuid.New(),
		Type:  api.PostProductsJSONBodyTypeОбувь,
	}

	tests := []test{
		{
			name:   "ok",
			input:  input,
			status: http.StatusCreated,
			res: &api.Product{
				DateTime:    &product.DateTime,
				Id:          &product.ID,
				ReceptionId: product.ReceptionID,
				Type:        api.ProductType(product.Type),
			},
			mockFn: func(test test, m mocks) {
				m.product.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&product, nil)
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
			name:   "forbidden",
			input:  input,
			status: http.StatusForbidden,
			res:    api.Error{Message: domain.ErrEmployeeOnly.Error()},
			mockFn: func(test test, m mocks) {
				m.product.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, domain.ErrEmployeeOnly)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "create product error",
			input:  input,
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.product.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errOops)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "no reception is progress error",
			input:  input,
			status: http.StatusBadRequest,
			res:    api.Error{Message: domain.ErrThereIsNoInProgressReception.Error()},
			mockFn: func(test test, m mocks) {
				m.product.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, domain.ErrThereIsNoInProgressReception)
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

			err := h.PostProducts(c)
			require.NoError(t, err)

			actual := rec.Body.String()
			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.JSONEq(t, string(expected)+"\n", actual)
			require.Equal(t, tt.status, rec.Code)
		})
	}
}

func TestProduct_PostPvzPvzIdDeleteLastProduct(t *testing.T) {
	t.Parallel()

	type test struct {
		name  string
		input openapi_types.UUID

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
			res:    nil,
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
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
			name:   "forbidden",
			input:  input,
			status: http.StatusForbidden,
			res:    api.Error{Message: domain.ErrEmployeeOnly.Error()},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain.ErrEmployeeOnly)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "create product error",
			input:  input,
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errOops)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "product not found",
			input:  input,
			status: http.StatusBadRequest,
			res:    api.Error{Message: domain.ErrProductNotFound.Error()},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain.ErrProductNotFound)
			},
			role: domain.UserRoleModerator,
		},
		{
			name:   "not found reception with status in progress ",
			input:  input,
			status: http.StatusBadRequest,
			res:    api.Error{Message: domain.ErrThereIsNoInProgressReception.Error()},
			mockFn: func(test test, m mocks) {
				m.reception.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain.ErrThereIsNoInProgressReception)
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

			err := h.PostPvzPvzIdDeleteLastProduct(c, tt.input)
			require.NoError(t, err)

			actual := rec.Body.String()
			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.Equal(t, tt.status, rec.Code)
			if tt.res == nil && actual == "" {
				return
			}
			require.JSONEq(t, string(expected)+"\n", actual)
		})
	}
}
