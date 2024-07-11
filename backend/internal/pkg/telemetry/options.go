package telemetry

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
