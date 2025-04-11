//go:generate ${LOCAL_BIN}/ifacemaker -f ./auth.go -s UseCase -i Auth -p mock_auth -o ./mocks/auth.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/auth.go -destination=./mocks/auth.go -package=mock_auth

package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	"github.com/khostya/pvz/pkg/hash"
)

type UseCase struct {
	userRepo       userRepo
	jwtManager     jwtManager
	passwordHasher passwordHasher
}

type (
	userRepo interface {
		Create(ctx context.Context, user *domain.User) (*domain.User, error)
		GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
		GetByEmail(ctx context.Context, email string) (*domain.User, error)
	}

	jwtManager interface {
		GenerateDummyToken(role domain.Role) (domain.Token, error)
		GenerateToken(user *domain.User) (domain.Token, error)
	}

	passwordHasher interface {
		Hash(password string) (string, error)
		Equal(param hash.EqualsParam) bool
	}

	AuthDepsUseCase struct {
		UserRepo       userRepo
		JwtManager     jwtManager
		PasswordHasher passwordHasher
	}
)

func NewAuthUseCase(useCase AuthDepsUseCase) *UseCase {
	return &UseCase{
		userRepo:       useCase.UserRepo,
		jwtManager:     useCase.JwtManager,
		passwordHasher: useCase.PasswordHasher,
	}
}

func (a *UseCase) Login(ctx context.Context, param dto.LoginUserParam) (domain.Token, error) {
	user, err := a.userRepo.GetByEmail(ctx, param.Email)
	if err != nil {
		return "", err
	}

	isEqual := a.passwordHasher.Equal(hash.EqualsParam{Hashed: user.Password, V: param.Password})
	if !isEqual {
		return "", domain.ErrInvalidPassword
	}

	return a.jwtManager.GenerateToken(user)
}

func (a *UseCase) Register(ctx context.Context, param dto.RegisterUserParam) (*domain.User, error) {
	password, err := a.passwordHasher.Hash(param.Password)
	if err != nil {
		return nil, err
	}

	return a.userRepo.Create(ctx, &domain.User{
		ID:       uuid.New(),
		Email:    param.Email,
		Password: password,
		Role:     param.Role,
	})
}

func (a *UseCase) DummyLogin(ctx context.Context, param dto.DummyLoginUserParam) (domain.Token, error) {
	return a.jwtManager.GenerateDummyToken(param.Role)
}
