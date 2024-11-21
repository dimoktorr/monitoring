package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type db interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

type transactionKey struct{}

type Repository struct {
	pgx     *pgxpool.Pool
	scanAPI *pgxscan.API
	builder squirrel.StatementBuilderType
}

func NewRepository(conn *pgxpool.Pool, scanApi *pgxscan.API) *Repository {
	return &Repository{
		pgx:     conn,
		scanAPI: scanApi,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *Repository) conn(ctx context.Context) db {
	if tx, ok := ctx.Value(transactionKey{}).(pgx.Tx); ok {
		return tx
	}

	return r.pgx
}

func (r *Repository) Exec(ctx context.Context, query string, args []interface{}) error {
	commandTag, err := r.conn(ctx).Exec(ctx, query, args...)
	log.Default().Print("query: %s, number of rows affected: %d", query, commandTag.RowsAffected())

	return err
}

func (r *Repository) QueryRows(ctx context.Context, dst interface{}, query string, args []interface{}) error {
	rows, err := r.conn(ctx).Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		return err
	}

	log.Default().Print("query rows: %s, rows affected: %d", query, rows.CommandTag().RowsAffected())
	defer rows.Close()

	err = r.scanAPI.ScanAll(dst, rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}

func (r *Repository) QueryRow(ctx context.Context, dst interface{}, query string, args []interface{}) error {
	row := r.conn(ctx).QueryRow(ctx, query, args...)
	log.Default().Print("query row: %s", query)

	return row.Scan(dst)
}

func (r *Repository) QueryOne(ctx context.Context, dst interface{}, query string, args []interface{}) error {
	rows, err := r.conn(ctx).Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		return err
	}
	log.Default().Print("query one: %s, rows affected: %d", query, rows.CommandTag().RowsAffected())
	defer rows.Close()

	err = r.scanAPI.ScanOne(dst, rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}

func startTracerSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer("postgres").Start(ctx, "repository."+spanName)
}
