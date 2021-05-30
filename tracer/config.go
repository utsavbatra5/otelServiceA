package tracer

import (
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Config specifies parameters relevant for otel trace provider.
type Config struct {
	// ApplicationName is the name for this application.
	ApplicationName string `json:"applicationName"`

	// Provider is the name of the trace provider to use.
	Provider string `json:"provider"`

	// Endpoint is the endpoint to which spans need to be submitted.
	Endpoint string `json:"endpoint"`

	// SkipTraceExport works only in case of provider stdout. Set
	// SkipTraceExport = true if you don't want to print the span
	// and tracer information in stdout.
	SkipTraceExport bool `json:"skipTraceExport"`

	// Providers are useful when client wants to add their own custom
	// TracerProvider.
	Providers map[string]ProviderConstructor `json:"-"`
}

// TraceConfig will be used in TraceMiddleware to use config and TraceProvider
// objects created by ConfigureTracerProvider.
// (Deprecated). Consider using Tracing instead.
type TraceConfig struct {
	TraceProvider trace.TracerProvider
}

// Tracing contains the core dependencies to make tracing possible across an
// application.
type Tracing struct {
	// Enabled should be set to false if the tracerProvider is a noop which
	// essentially disables tracing in the system.
	Enabled bool

	// TracerProvider helps create trace spans.
	TracerProvider trace.TracerProvider

	// Propagator helps propagate trace context across API boundaries.
	Propagator propagation.TextMapPropagator
}
