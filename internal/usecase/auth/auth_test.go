package auth

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	mock_jwt "github.com/khostya/pvz/internal/lib/jwt/mocks"
	mock_postgres "github.com/khostya/pvz/internal/repo/postgres/mocks"
	"github.com/khostya/pvz/pkg/hash"
	mock_hash "github.com/khostya/pvz/pkg/hash/mocks"
	mock_transactor "github.com/khostya/pvz/pkg/postgres/transactor/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	ErrAuthOops = errors.New("auth oops error")
)

type authMocks struct {
	mockUserRepo   *mock_postgres.MockUserRepo
	mockBcrypt     *mock_hash.Mockbcrypt
	mockJWT        *mock_jwt.Mockmanager
	mockTransactor *mock_transactor.MockTransactor
}

func newAuthMocks(t *testing.T) authMocks {
	ctrl := gomock.NewController(t)

	return authMocks{
		mockTransactor: mock_transactor.NewMockTransactor(ctrl),
		mockUserRepo:   mock_postgres.NewMockUserRepo(ctrl),
		mockBcrypt:     mock_hash.NewMockbcrypt(ctrl),
		mockJWT:        mock_jwt.NewMockmanager(ctrl),
	}
}

func newDepsAuthUseCase(mocks authMocks) AuthDepsUseCase {
	return AuthDepsUseCase{
		UserRepo:       mocks.mockUserRepo,
		PasswordHasher: mocks.mockBcrypt,
		JwtManager:     mocks.mockJWT,
	}
}

func TestAuthUseCase_DummyLogin(t *testing.T) {
	t.Parallel()

	type test struct {
		name    string
		input   dto.DummyLoginUserParam
		token   domain.Token
		mockFn  func(ctx context.Context, test test, m authMocks)
		wantErr error
	}

	ctx := context.Background()
	tests := []test{
		{
			name:    "employee is ok",
			input:   dto.DummyLoginUserParam{Role: domain.UserRoleEmployee},
			wantErr: nil,
			token:   domain.Token(uuid.New().String()),
			mockFn: func(ctx context.Context, test test, m authMocks) {
				m.mockJWT.EXPECT().GenerateDummyToken(test.input.Role).
					Return(test.token, nil)
			},
		},
		{
			name:    "moderator is ok",
			input:   dto.DummyLoginUserParam{Role: domain.UserRoleModerator},
			wantErr: nil,
			token:   domain.Token(uuid.New().String()),
			mockFn: func(ctx context.Context, test test, m authMocks) {
				m.mockJWT.EXPECT().GenerateDummyToken(test.input.Role).
					Return(test.token, nil).
					Times(1)
			},
		},
		{
			name:    "employee is error",
			input:   dto.DummyLoginUserParam{Role: domain.UserRoleEmployee},
			wantErr: ErrAuthOops,
			token:   "",
			mockFn: func(ctx context.Context, test test, m authMocks) {
				m.mockJWT.EXPECT().GenerateDummyToken(test.input.Role).
					Return(domain.Token(""), ErrAuthOops).
					Times(1)
			},
		},
		{
			name:    "moderator is error",
			input:   dto.DummyLoginUserParam{Role: domain.UserRoleModerator},
			wantErr: ErrAuthOops,
			token:   "",
			mockFn: func(ctx context.Context, test test, m authMocks) {
				m.mockJWT.EXPECT().GenerateDummyToken(test.input.Role).
					Return(domain.Token(""), ErrAuthOops).
					Times(1)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newAuthMocks(t)

			authUseCase := NewAuthUseCase(newDepsAuthUseCase(mocks))

			tt.mockFn(ctx, tt, mocks)

			token, err := authUseCase.DummyLogin(ctx, tt.input)
			require.Equal(t, tt.wantErr, err)

			require.Equal(t, tt.token, token)
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()
	)

	type test struct {
		name           string
		input          dto.LoginUserParam
		hashedPassword string
		token          domain.Token
		err            error
		mockFn         func(ctx context.Context, test test, m authMocks)
	}

	input := dto.LoginUserParam{
		Email:    gofakeit.Email(),
		Password: uuid.New().String(),
	}
	tests := []test{
		{
			name:           "error get user by email",
			input:          input,
			token:          "",
			hashedPassword: uuid.New().String(),
			mockFn: func(ctx context.Context, test test, m authMocks) {
				m.mockUserRepo.EXPECT().GetByEmail(ctx, test.input.Email).
					Times(1).Return(nil, ErrAuthOops)
			},
			err: ErrAuthOops,
		},
		{
			name:           "incorrect password",
			input:          input,
			token:          "",
			hashedPassword: uuid.New().String(),
			mockFn: func(ctx context.Context, test test, m authMocks) {
				user := &domain.User{Password: test.hashedPassword}

				m.mockUserRepo.EXPECT().GetByEmail(ctx, test.input.Email).
					Times(1).Return(user, nil)
				m.mockBcrypt.EXPECT().Equal(hash.EqualsParam{Hashed: test.hashedPassword, V: test.input.Password}).
					Return(false).Times(1)
			},
			err: domain.ErrInvalidPassword,
		},
		{
			name:           "error generate token",
			input:          input,
			token:          "",
			hashedPassword: uuid.New().String(),
			mockFn: func(ctx context.Context, test test, m authMocks) {
				user := &domain.User{Password: test.hashedPassword}

				m.mockUserRepo.EXPECT().GetByEmail(ctx, test.input.Email).
					Times(1).Return(user, nil)
				m.mockBcrypt.EXPECT().Equal(hash.EqualsParam{Hashed: test.hashedPassword, V: test.input.Password}).
					Return(true).Times(1)
				m.mockJWT.EXPECT().GenerateToken(user).
					Times(1).Return(domain.Token(""), ErrAuthOops)
			},
			err: ErrAuthOops,
		},
		{
			name:           "ok",
			input:          input,
			token:          domain.Token(uuid.New().String()),
			hashedPassword: uuid.New().String(),
			mockFn: func(ctx context.Context, test test, m authMocks) {
				user := &domain.User{Password: test.hashedPassword}

				m.mockUserRepo.EXPECT().GetByEmail(ctx, test.input.Email).
					Times(1).Return(user, nil)
				m.mockBcrypt.EXPECT().Equal(hash.EqualsParam{Hashed: test.hashedPassword, V: test.input.Password}).
					Return(true).Times(1)
				m.mockJWT.EXPECT().GenerateToken(user).
					Return(test.token, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newAuthMocks(t)
			tt.mockFn(ctx, tt, mocks)

			authUseCase := NewAuthUseCase(newDepsAuthUseCase(mocks))

			token, err := authUseCase.Login(ctx, tt.input)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.token, token)
		})
	}
}

func TestAuthUseCase_Register(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()
	)

	type test struct {
		name   string
		input  dto.RegisterUserParam
		user   *domain.User
		err    error
		mockFn func(ctx context.Context, test test, m authMocks)
	}

	hashedPassword := uuid.New().String()
	input := dto.RegisterUserParam{
		Email:    "sgseegspoko",
		Password: "sgesg",
	}
	tests := []test{
		{
			name:  "ok",
			input: input,
			user:  &domain.User{ID: uuid.New(), Password: hashedPassword},
			mockFn: func(ctx context.Context, test test, m authMocks) {
				m.mockBcrypt.EXPECT().Hash(test.input.Password).
					Times(1).Return(hashedPassword, nil)
				m.mockUserRepo.EXPECT().Create(ctx, gomock.Any()).
					Times(1).Return(test.user, nil)
			},
		},
		{
			name:  "error hash",
			input: input,
			mockFn: func(ctx context.Context, test test, m authMocks) {
				m.mockBcrypt.EXPECT().Hash(test.input.Password).
					Times(1).Return("", ErrAuthOops)
			},
			err: ErrAuthOops,
		},
		{
			name:  "error create user",
			input: input,
			mockFn: func(ctx context.Context, test test, m authMocks) {
				m.mockBcrypt.EXPECT().Hash(test.input.Password).
					Times(1).Return(hashedPassword, nil)
				m.mockUserRepo.EXPECT().Create(ctx, gomock.Any()).
					Times(1).Return(test.user, ErrAuthOops)
			},
			err: ErrAuthOops,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newAuthMocks(t)
			tt.mockFn(ctx, tt, mocks)

			authUseCase := NewAuthUseCase(newDepsAuthUseCase(mocks))

			user, err := authUseCase.Register(ctx, tt.input)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.user, user)
		})
	}
}
