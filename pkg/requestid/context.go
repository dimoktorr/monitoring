package requestid

import (
	"context"
)

type ContextKey struct{}

func WithContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKey{}, requestID)
}

func FromContext(ctx context.Context) string {
	val := ctx.Value(ContextKey{})
	if val == nil {
		return ""
	}

	requestID, ok := val.(string)
	if !ok {
		return ""
	}

	return requestID
}
