package api

import (
	"context"
	"github.com/dimoktorr/monitoring/internal/app"
	v1 "github.com/dimoktorr/monitoring/pkg/api/v1"
	"github.com/dimoktorr/monitoring/pkg/requestid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	v1.UnimplementedExampleServiceServer
	service *app.Service
}

func NewServer(service *app.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Pay(ctx context.Context, in *v1.PayRequest) (*v1.PayResponse, error) {
	log.Println("PayProduct started", "request_id", requestid.FromContext(ctx))
	defer log.Println("PayProduct finished", "request_id", requestid.FromContext(ctx))

	ctx, span := startTracerSpan(ctx, "PayProduct")
	defer span.End()

	span.AddEvent("", trace.WithAttributes(
		attribute.Int("id_product", int(in.GetProductId())),
		attribute.String("request_id", requestid.FromContext(ctx)),
	))

	payStatus, err := s.service.PayProduct(ctx, int(in.GetProductId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "can't get product")
	}

	return &v1.PayResponse{
		Status: payStatus,
	}, nil
}

func (s *Server) GetProduct(ctx context.Context, in *v1.GetRequest) (*v1.GetResponse, error) {
	log.Println("GetProduct started", "request_id", requestid.FromContext(ctx))
	defer log.Println("GetProduct finished", "request_id", requestid.FromContext(ctx))

	ctx, span := startTracerSpan(ctx, "GetProduct")
	defer span.End()

	span.AddEvent("", trace.WithAttributes(
		attribute.Int("id_product", int(in.GetId())),
		attribute.String("request_id", requestid.FromContext(ctx)),
	))

	product, err := s.service.GetProduct(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "can't get product")
	}

	return &v1.GetResponse{
		Products: []*v1.Product{
			{
				Id:    int32(product.ID),
				Name:  product.Name,
				Price: float32(product.Price),
			},
		},
	}, nil
}

func startTracerSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer("api").Start(ctx, "apiV1."+spanName)
}
