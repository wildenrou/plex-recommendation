package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	spanTracer "go.opentelemetry.io/otel/trace"
	"time"

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
	withName    string
	withPackage string
}

type SpanOption func(*spanOption)

func WithSpanName(s string) SpanOption {
	return func(o *spanOption) {
		o.withName = s
	}
}

func WithSpanPackage(s string) SpanOption {
	return func(o *spanOption) {
		o.withPackage = s
	}
}

func StartSpan(ctx context.Context, opts ...SpanOption) (context.Context, spanTracer.Span) {
	opt := &spanOption{}
	for _, o := range opts {
		o(opt)
	}
	spanName := "Span Name Unspecified"
	if opt.withName != "" {
		spanName = opt.withName
	}

	packageName := "Unspecified"
	if opt.withPackage != "" {
		packageName = opt.withPackage
	}
	ctx, span := tracer.Start(ctx, spanName)
	span.SetAttributes(attribute.String("package", packageName))
	return ctx, span
}
