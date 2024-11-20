package app

import (
	"context"
	"github.com/dimoktorr/monitoring/pkg/requestid"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryRequestIDServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		meta, isSet := metadata.FromIncomingContext(ctx)
		if !isSet {
			return handler(requestid.WithContext(ctx, uuid.New().String()), req)
		}

		requestID := requestid.FromGRPCMetadata(meta)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		return handler(requestid.WithContext(ctx, requestID), req)
	}
}
