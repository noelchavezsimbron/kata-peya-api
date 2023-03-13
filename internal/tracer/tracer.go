package tracer

import (
	"context"

	"kata-peya/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	conf = config.Get()
	t    = otel.Tracer(conf.InstanceName)
)

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.Start(ctx, spanName, opts...)
}
