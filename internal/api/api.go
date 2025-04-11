package api

import (
	"context"
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/khostya/pvz/internal/api/v1/grpc"
	"github.com/khostya/pvz/internal/api/v1/http"
	middleware2 "github.com/khostya/pvz/internal/api/v1/http/middleware"
	api "github.com/khostya/pvz/internal/api/v1/http/server"
	cache2 "github.com/khostya/pvz/internal/cache"
	"github.com/khostya/pvz/internal/config"
	"github.com/khostya/pvz/internal/lib/jwt"
	"github.com/khostya/pvz/internal/usecase"
	pvz_v1 "github.com/khostya/pvz/pkg/api/v1/proto"
	grpcserver "github.com/khostya/pvz/pkg/grpc"
	"github.com/khostya/pvz/pkg/httpserver"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	h "net/http"
)

type (
	deps struct {
		Manager *jwt.Manager
		uc      *usecase.UseCase
	}
)

func New(ctx context.Context, cfg config.Config, uc *usecase.UseCase, manager *jwt.Manager) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	deps := deps{
		Manager: manager,
		uc:      uc,
	}

	srv := newHttpServer(cfg.App, cfg, deps)
	srv.Start()

	grpcserver, err := newGRPCServer(ctx, cfg, deps)
	if err != nil {
		return err
	}

	err = grpcserver.Start()
	if err != nil {
		return err
	}

	promSrv := newPrometheusMetrics(cfg.API.Prometheus)
	promSrv.Start()

	select {
	case <-ctx.Done():
	case srvErr := <-srv.Notify():
		err = multierror.Append(err, srvErr)
	case srvErr := <-promSrv.Notify():
		err = multierror.Append(err, srvErr)
	case srvErr := <-grpcserver.Wait():
		err = multierror.Append(err, srvErr)
	}
	cancel()

	if srvErr := srv.Shutdown(ctx); srvErr != nil && !errors.Is(srvErr, h.ErrServerClosed) {
		err = multierror.Append(err, srvErr)
	}

	if srvErr := promSrv.Shutdown(ctx); srvErr != nil && !errors.Is(srvErr, h.ErrServerClosed) {
		err = multierror.Append(err, srvErr)
	}

	if srvErr := <-grpcserver.Wait(); srvErr != nil {
		err = multierror.Append(err, srvErr)
	}

	return err
}

func newHttpServer(app config.App, cfg config.Config, deps deps) *httpserver.Server {
	cfgHttp := cfg.API.HTTP

	v1 := echo.New()

	server := http.NewServer(http.Deps{
		Product:   deps.uc.Product,
		Pvz:       deps.uc.Pvz,
		Reception: deps.uc.Reception,
		Auth:      deps.uc.Auth,
	})

	v1.Use(echoprometheus.NewMiddleware(app.Name))
	v1.Use(middleware.Recover())
	v1.Use(middleware.Logger())
	v1.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	srv := httpserver.New(v1,
		httpserver.Port(cfgHttp.Port),
		httpserver.ReadTimeout(cfgHttp.ReadTimeout),
		httpserver.WriteTimeout(cfgHttp.WriteTimeout),
		httpserver.IdleTimeout(cfgHttp.IdleTimeout),
	)

	mw, err := middleware2.CreateValidatorMiddleware(middleware2.NewAuthenticator(deps.Manager))
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}

	v1.Use(mw)

	api.RegisterHandlers(v1, server)

	return srv
}

func newPrometheusMetrics(cfg config.Prometheus) *httpserver.Server {
	v1 := echo.New()

	v1.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	v1.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	srv := httpserver.New(v1, httpserver.Port(cfg.Port))

	return srv
}

func newGRPCServer(ctx context.Context, cfg config.Config, deps deps) (*grpcserver.Server, error) {
	getPVZListResponse, err := cache2.NewPvzList[string, *pvz_v1.GetPVZListResponse](cfg.Cache.PvzListTTL)
	if err != nil {
		return nil, err
	}

	grpcserver := grpcserver.New(ctx, int(cfg.API.GRPC.Port))

	grpcService := grpc.NewServer(grpc.Deps{
		PvzService:         deps.uc.Pvz,
		GetPVZListResponse: getPVZListResponse,
	})
	grpcService.Register(grpcserver.GetRegistrar())

	return grpcserver, nil
}
