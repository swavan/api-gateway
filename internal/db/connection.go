package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Api interface {
	Close(ctx context.Context)
	SetPool(p *pgxpool.Pool)
	GetPool() *pgxpool.Pool
	GetDB() *sqlx.DB
}

type database struct {
	pool *pgxpool.Pool
}

func Start(ctx context.Context, url string) (Api, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	return &database{pool}, err
}

func (d *database) Close(ctx context.Context) {
	d.pool.Close()
}

func (d *database) SetPool(p *pgxpool.Pool) {
	d.pool = p
}

func (d *database) GetPool() *pgxpool.Pool {
	return d.pool
}

func (d *database) GetDB() *sqlx.DB {
	connector := stdlib.OpenDBFromPool(d.pool)
	return sqlx.NewDb(connector, "pgx")
}
