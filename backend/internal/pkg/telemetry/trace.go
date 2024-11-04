package telemetry

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	spanTracer "go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tracer = otel.Tracer(otelName)

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider() (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("Plex Recommendation Backend"),
			),
		),
	)
	return traceProvider, nil
}

type spanOption struct {
	name        string
	packageName string
	requestId   string
}

type SpanOption func(*spanOption)

func WithSpanName(s string) SpanOption {
	return func(o *spanOption) {
		o.name = s
	}
}

func WithSpanPackage(s string) SpanOption {
	return func(o *spanOption) {
		o.packageName = s
	}
}

func WithRequestId(s string) SpanOption {
	return func(o *spanOption) {
		o.requestId = s
	}
}

// StartSpan will initialize a new telemetry span with the provided context
// and options. If telemetry is explicity disabled via DISABLE_TELEMETRY,
// it will return the same context back with a no-op span.
func StartSpan(ctx context.Context, opts ...SpanOption) (context.Context, spanTracer.Span) {
	if noOp {
		return ctx, spanTracer.SpanFromContext(ctx)
	}
	opt := &spanOption{}
	for _, o := range opts {
		o(opt)
	}
	spanName := "Span Name Unspecified"
	if opt.name != "" {
		spanName = opt.name
	}

	packageName := "Unspecified"
	if opt.packageName != "" {
		packageName = opt.packageName
	}
	requestId := uuid.NewString()
	if opt.requestId != "" {
		requestId = opt.requestId
	}
	ctx, span := tracer.Start(ctx, spanName)
	span.SetAttributes(attribute.String("package", packageName))
	span.SetAttributes(attribute.String("request_id", requestId))
	return ctx, span
}
