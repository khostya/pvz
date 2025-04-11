package api

import (
	"context"
	"github.com/khostya/pvz/internal/api/v1/grpc"
	"github.com/khostya/pvz/internal/api/v1/http"
	middleware2 "github.com/khostya/pvz/internal/api/v1/http/middleware"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	"github.com/khostya/pvz/internal/config"
	"github.com/khostya/pvz/internal/lib/jwt"
	"github.com/khostya/pvz/internal/usecase"
	grpcserver "github.com/khostya/pvz/pkg/grpc"
	"github.com/khostya/pvz/pkg/httpserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

type (
	deps struct {
		Manager *jwt.Manager
		uc      *usecase.UseCase
	}
)

func New(ctx context.Context, cfg config.API, uc *usecase.UseCase, manager *jwt.Manager) error {
	deps := deps{
		Manager: manager,
		uc:      uc,
	}

	srv := newHttpServer(cfg.HTTP, deps)
	srv.Start()

	grpcserver := newGRPCServer(ctx, cfg.GRPC, deps)
	err := grpcserver.Start()
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	case <-srv.Notify():
	case <-grpcserver.Wait():
	}

	return nil
}

func newHttpServer(cfg config.HTTP, deps deps) *httpserver.Server {
	v1 := echo.New()

	server := http.NewServer(http.Deps{
		Product:   deps.uc.Product,
		Pvz:       deps.uc.Pvz,
		Reception: deps.uc.Reception,
		Auth:      deps.uc.Auth,
	})

	v1.Use(middleware.Recover())
	v1.Use(middleware.Logger())
	v1.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	srv := httpserver.New(v1,
		httpserver.Port(cfg.Port),
		httpserver.ReadTimeout(cfg.ReadTimeout),
		httpserver.WriteTimeout(cfg.WriteTimeout))

	mw, err := middleware2.CreateValidatorMiddleware(middleware2.NewAuthenticator(deps.Manager))
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}

	v1.Use(mw)

	api.RegisterHandlers(v1, server)

	return srv
}

func newGRPCServer(ctx context.Context, cfg config.GRPC, deps deps) *grpcserver.Server {
	grpcserver := grpcserver.New(ctx, int(cfg.Port))
	grpcService := grpc.NewServer(deps.uc.Pvz)
	grpcService.Register(grpcserver.GetRegistrar())

	return grpcserver
}
