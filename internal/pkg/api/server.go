package api

import (
	"context"
	"github.com/dimoktorr/monitoring/internal/app"
	v1 "github.com/dimoktorr/monitoring/pkg/api/v1"
	"github.com/dimoktorr/monitoring/pkg/requestid"
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

func (s *Server) GetProduct(ctx context.Context, in *v1.GetRequest) (*v1.GetResponse, error) {
	log.Println("GetProduct started", "request_id", requestid.FromContext(ctx))
	defer log.Println("GetProduct finished", "request_id", requestid.FromContext(ctx))

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
