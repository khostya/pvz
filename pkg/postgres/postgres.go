package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khostya/pvz/internal/config"
	"github.com/pkg/errors"
)

type Pool = pgxpool.Pool

func NewPool(ctx context.Context, pgCfg config.PG) (*Pool, error) {
	cfg, err := pgxpool.ParseConfig(pgCfg.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse connection URL")
	}

	cfg.MaxConnLifetime = pgCfg.ConnMaxLifetime
	cfg.MaxConnIdleTime = pgCfg.ConnMaxIdleTime
	cfg.MaxConns = pgCfg.MaxOpenConns
	cfg.MinIdleConns = pgCfg.MinIdleConns

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create connection pool")
	}
	return pool, nil
}

func NewPoolFromURL(ctx context.Context, url string) (*Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create connection pool")
	}
	return pool, nil
}
