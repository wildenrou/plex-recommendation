package telemetry

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

const otelName = "Plex Recommendation API"

var (
	Tracer = otel.Tracer(otelName)
	Meter  = otel.Meter(otelName)
)

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
