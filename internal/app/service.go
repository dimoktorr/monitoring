package app

import (
	"context"
	"github.com/dimoktorr/monitoring/internal/domain"
	"github.com/dimoktorr/monitoring/pkg/requestid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type Service struct {
	metrics Metrics
	repo    Repository
	storage Storage
}

func New(
	metrics Metrics,
	repo Repository,
	storage Storage,
) *Service {
	return &Service{
		metrics: metrics,
		repo:    repo,
		storage: storage,
	}
}

func (s *Service) GetProduct(ctx context.Context, productId int) (*domain.Product, error) {
	ctx, span := startTracerSpan(ctx, "GetProduct")
	defer span.End()

	span.AddEvent("Service", trace.WithAttributes(
		attribute.Int("id_product", productId),
		attribute.String("request_id", requestid.FromContext(ctx)),
	))

	product, err := s.repo.GetProduct(ctx, productId)
	if err != nil {
		log.Default().Println("error get product", err, "request_id", requestid.FromContext(ctx))

		span.RecordError(err, trace.WithAttributes(
			attribute.Int("id_product", productId),
			attribute.String("request_id", requestid.FromContext(ctx)),
		))

		return nil, err
	}

	s.metrics.IncGetProductSumCounter()

	return product, nil
}

func startTracerSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer("apiExampleV1").Start(ctx, "apiV1."+spanName)
}
