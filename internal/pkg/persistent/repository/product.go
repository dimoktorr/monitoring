package repository

import (
	"context"
	"github.com/dimoktorr/monitoring/internal/domain"
	"github.com/dimoktorr/monitoring/pkg/requestid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (r *Repository) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	ctx, span := startTracerSpan(ctx, "GetProduct")
	defer span.End()

	span.AddEvent("", trace.WithAttributes(
		attribute.Int("id_product", id),
		attribute.String("request_id", requestid.FromContext(ctx)),
	))

	query, args, err := r.builder.
		Select("id", "name", "price").
		From("products").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return nil, err
	}

	var product Product

	if err := r.QueryOne(ctx, &product, query, args); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.Int("id_product", id),
			attribute.String("request_id", requestid.FromContext(ctx)),
		))

		return nil, err
	}
	return &domain.Product{
		ID:    int(product.ID.Int),
		Name:  product.Name.String,
		Price: product.Price.Float,
	}, nil
}
