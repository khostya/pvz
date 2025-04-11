package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	mock_cache "github.com/khostya/pvz/internal/cache/mocks"
	"github.com/khostya/pvz/internal/domain"
	mock_auth "github.com/khostya/pvz/internal/usecase/auth/mocks"
	mock_product "github.com/khostya/pvz/internal/usecase/product/mocks"
	mock_pvz "github.com/khostya/pvz/internal/usecase/pvz/mocks"
	mock_reception "github.com/khostya/pvz/internal/usecase/reception/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	errOops = errors.New("oops error")
	e       = echo.New()
)

type mocks struct {
	product             *mock_product.MockProduct
	auth                *mock_auth.MockAuth
	pvz                 *mock_pvz.MockPvz
	reception           *mock_reception.MockReception
	getPvzResponseCache *mock_cache.MockCache[string, []getPvzResponse]
}

func newMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)
	return mocks{
		product:             mock_product.NewMockProduct(ctrl),
		auth:                mock_auth.NewMockAuth(ctrl),
		pvz:                 mock_pvz.NewMockPvz(ctrl),
		reception:           mock_reception.NewMockReception(ctrl),
		getPvzResponseCache: mock_cache.NewMockCache[string, []getPvzResponse](ctrl),
	}
}

func NewMockServer(m mocks) *Server {
	return NewServer(Deps{
		Reception:           m.reception,
		Auth:                m.auth,
		Product:             m.product,
		Pvz:                 m.pvz,
		GetPvzResponseCache: m.getPvzResponseCache,
	})
}

func TestAuth_PostDummyLogin(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  *api.PostDummyLoginJSONBody
		status int
		res    any
		mockFn func(test test, m mocks)
		token  string
	}

	input := &api.PostDummyLoginJSONBody{
		Role: api.PostDummyLoginJSONBodyRole(api.Moderator),
	}
	token := uuid.New().String()

	tests := []test{
		{
			name:   "ok",
			input:  input,
			status: http.StatusOK,
			res:    token,
			token:  token,
			mockFn: func(test test, m mocks) {
				m.auth.EXPECT().DummyLogin(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain.Token(test.token), nil)
			},
		},
		{
			name:   "dummy login error",
			input:  input,
			token:  "",
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.auth.EXPECT().DummyLogin(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain.Token(""), errOops)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, h, rec, mocks := newEchoCtx(t, tt.input)
			tt.mockFn(tt, mocks)

			err := h.PostDummyLogin(c)
			require.NoError(t, err)

			actual := rec.Body.String()

			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.Equal(t, tt.status, rec.Code)
			require.Equal(t, string(expected)+"\n", actual)
		})
	}
}

func newEchoCtx(t *testing.T, input any) (echo.Context, *Server, *httptest.ResponseRecorder, mocks) {
	body, err := json.Marshal(input)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mocks := newMocks(t)

	h := NewMockServer(mocks)

	return c, h, rec, mocks
}

func TestAuth_PostLogin(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  *api.PostLoginJSONBody
		status int
		res    any
		mockFn func(test test, m mocks)
		token  string
	}

	input := &api.PostLoginJSONBody{
		Email:    "khostya.konsantin@gmail.com",
		Password: "bla-bla-bla",
	}
	token := uuid.New().String()

	tests := []test{
		{
			name:   "ok",
			input:  input,
			status: http.StatusOK,
			res:    token,
			token:  token,
			mockFn: func(test test, m mocks) {
				m.auth.EXPECT().Login(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain.Token(test.token), nil)
			},
		},
		{
			name:   "login error",
			input:  input,
			token:  "",
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.auth.EXPECT().Login(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain.Token(""), errOops)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, h, rec, mocks := newEchoCtx(t, tt.input)
			tt.mockFn(tt, mocks)

			err := h.PostLogin(c)
			require.NoError(t, err)

			actual := rec.Body.String()

			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.Equal(t, tt.status, rec.Code)
			require.Equal(t, string(expected)+"\n", actual)
		})
	}
}

func TestAuth_PostRegister(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  *api.PostRegisterJSONBody
		status int
		res    any
		mockFn func(test test, m mocks)
	}

	input := &api.PostRegisterJSONBody{
		Email:    "khostya.konsantin@gmail.com",
		Password: "bla-bla-bla",
		Role:     api.Moderator,
	}

	tests := []test{
		{
			name:   "ok",
			input:  input,
			status: http.StatusCreated,
			res: &api.User{
				Email: "khostya.konsantin@gmail.com",
				Id:    &user.ID,
				Role:  api.UserRole(user.Role),
			},
			mockFn: func(test test, m mocks) {
				m.auth.EXPECT().Register(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
		},
		{
			name:   "register error",
			input:  input,
			status: http.StatusInternalServerError,
			res:    api.Error{Message: errOops.Error()},
			mockFn: func(test test, m mocks) {
				m.auth.EXPECT().Register(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errOops)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, h, rec, mocks := newEchoCtx(t, tt.input)
			tt.mockFn(tt, mocks)

			err := h.PostRegister(c)
			require.NoError(t, err)

			actual := rec.Body.String()
			expected, err := json.Marshal(tt.res)
			require.NoError(t, err)

			require.JSONEq(t, string(expected)+"\n", actual)
			require.Equal(t, tt.status, rec.Code)
		})
	}
}
