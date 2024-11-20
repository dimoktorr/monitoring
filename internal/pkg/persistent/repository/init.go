package repository

import (
	"context"
	"database/sql"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConn struct {
	Pgx     *pgxpool.Pool
	ScanAPI *pgxscan.API
}

func NewPostgresConn(ctx context.Context, cfg Config) (*PostgresConn, error) {
	pgxConn, err := New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	scanAPI, err := newScanApi()
	if err != nil {
		return nil, err
	}

	return &PostgresConn{
		Pgx:     pgxConn,
		ScanAPI: scanAPI,
	}, nil
}

func newScanApi() (*pgxscan.API, error) {
	scanner, err := pgxscan.NewDBScanAPI(
		dbscan.WithScannableTypes((*sql.Scanner)(nil)),
	)
	if err != nil {
		return nil, err
	}

	return pgxscan.NewAPI(scanner)
}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	c, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	c.MaxConnIdleTime = cfg.MaxConnIdleTime

	conn, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
