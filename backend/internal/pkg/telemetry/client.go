package telemetry

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const otelName = "Plex Recommendation API"

var (
	Tracer = otel.Tracer(otelName)
	Meter  = otel.Meter(otelName)
)

type telemetryOptions struct {
	withTracer bool
	withMeter  bool
}

type Option func(*telemetryOptions)

func WithTracer(enabled bool) Option {
	return func(o *telemetryOptions) {
		o.withTracer = enabled
	}
}

func WithMeter(enabled bool) Option {
	return func(o *telemetryOptions) {
		o.withMeter = enabled
	}
}

// InitOtel bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func InitOtel(ctx context.Context, opts ...Option) (func(context.Context) error, error) {

	var options = telemetryOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	var err error
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	if options.withTracer {
		var tracerProvider *trace.TracerProvider
		tracerProvider, err = newTraceProvider()
		if err != nil {
			handleErr(err)
			return shutdown, err
		}
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
	}

	// Set up meter provider.
	if options.withMeter {
		var meterProvider *metric.MeterProvider
		meterProvider, err = newMeterProvider()
		if err != nil {
			handleErr(err)
			return shutdown, err
		}
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
	}

	return shutdown, nil
}

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

func newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)))
	return meterProvider, nil
}
