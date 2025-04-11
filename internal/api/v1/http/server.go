package http

import (
	"context"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
)

type (
	Server struct {
		product   productService
		auth      authService
		reception receptionService
		pvz       pvzService
	}

	productService interface {
		Create(ctx context.Context, param dto.CreateProductParam) (*domain.Product, error)
	}

	authService interface {
		Login(ctx context.Context, param dto.LoginUserParam) (domain.Token, error)
		DummyLogin(ctx context.Context, param dto.DummyLoginUserParam) (domain.Token, error)
		Register(ctx context.Context, param dto.RegisterUserParam) (*domain.User, error)
	}

	pvzService interface {
		GetPvz(ctx context.Context, param dto.GetPvzParam) ([]*domain.PVZ, error)
		Create(ctx context.Context, param dto.CreatePvzParam) (*domain.PVZ, error)
	}

	receptionService interface {
		DeleteLastProduct(ctx context.Context, param dto.DeleteLastReceptionParam) error
		CloseLastReception(ctx context.Context, param dto.CloseLastReceptionParam) (*domain.Reception, error)
		Create(ctx context.Context, param dto.CreateReceptionParam) (*domain.Reception, error)
	}
)

type Deps struct {
	Product   productService
	Auth      authService
	Pvz       pvzService
	Reception receptionService
}

func NewServer(deps Deps) *Server {
	return &Server{
		product:   deps.Product,
		auth:      deps.Auth,
		pvz:       deps.Pvz,
		reception: deps.Reception,
	}
}
