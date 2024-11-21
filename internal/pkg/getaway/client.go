package getaway

import (
	"context"
	"github.com/dimoktorr/monitoring/pkg/requestid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Getaway struct {
	metrics Metrics
}

func NewGetaway(metrics Metrics) *Getaway {
	return &Getaway{
		metrics: metrics,
	}
}

func (g *Getaway) Pay(ctx context.Context, amount float64) (string, error) {
	ctx, span := startTracerSpan(ctx, "Pay")
	defer span.End()

	span.AddEvent("", trace.WithAttributes(
		attribute.String("request_id", requestid.FromContext(ctx)),
	))

	//TODO: implement payment logic
	status := "success"

	g.metrics.IncPayProductSumCounter(status)

	if status == "success" {
		g.metrics.AddAmountPayProduct(amount)
	}

	return status, nil
}

func startTracerSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer("gatewayClient").Start(ctx, "gatewayV1."+spanName)
}
