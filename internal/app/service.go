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
	gateway Gateway
}

func New(
	metrics Metrics,
	repo Repository,
	storage Storage,
	gateway Gateway,
) *Service {
	return &Service{
		metrics: metrics,
		repo:    repo,
		storage: storage,
		gateway: gateway,
	}
}

func (s *Service) GetProduct(ctx context.Context, productId int) (*domain.Product, error) {
	ctx, span := startTracerSpan(ctx, "GetProduct")
	defer span.End()

	span.AddEvent("", trace.WithAttributes(
		attribute.Int("id_product", productId),
		attribute.String("request_id", requestid.FromContext(ctx)),
	))

	product, err := s.repo.GetProduct(ctx, productId)
	if err != nil {
		log.Default().Println("error get product", err, "request_id", requestid.FromContext(ctx))
		return nil, err
	}

	s.metrics.IncGetProductSumCounter()

	return product, nil
}

func (s *Service) PayProduct(ctx context.Context, productId int) (string, error) {
	ctx, span := startTracerSpan(ctx, "PayProduct")
	defer span.End()

	span.AddEvent("", trace.WithAttributes(
		attribute.Int("id_product", productId),
		attribute.String("request_id", requestid.FromContext(ctx)),
	))

	product, err := s.repo.GetProduct(ctx, productId)
	if err != nil {
		log.Default().Println("error get product", err, "request_id", requestid.FromContext(ctx))
		return "", err
	}

	status, payErr := s.gateway.Pay(ctx, product.Price)
	if payErr != nil {
		return "", payErr
	}

	return status, nil
}

func startTracerSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer("service").Start(ctx, "serviceV1."+spanName)
}
