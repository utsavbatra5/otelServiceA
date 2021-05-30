package tracer

import (
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrTracerProviderNotFound    = errors.New("TracerProvider builder could not be found")
	ErrTracerProviderBuildFailed = errors.New("Failed building TracerProvider")
)

// DefaultTracerProvider is used when no provider is given.
// The Noop tracer provider turns all tracing related operations into
// noops essentially disabling tracing.
const DefaultTracerProvider = "noop"

// ConfigureTracerProvider creates the TracerProvider based on the configuration
// provided. It has built-in support for jaeger, zipkin, stdout and noop providers.
// A different provider can be used if a constructor for it is provided in the
// config.
// If a provider name is not provided, a noop tracerProvider will be returned.
func ConfigureTracerProvider(config Config) (trace.TracerProvider, error) {
	if len(config.Provider) == 0 {
		config.Provider = DefaultTracerProvider
	}
	// Handling camelcase of provider.
	config.Provider = strings.ToLower(config.Provider)
	providerConfig := config.Providers[config.Provider]
	if providerConfig == nil {
		providerConfig = providersConfig[config.Provider]
	}
	if providerConfig == nil {
		return nil, fmt.Errorf("%w for provider %s", ErrTracerProviderNotFound, config.Provider)
	}
	provider, err := providerConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTracerProviderBuildFailed, err)
	}
	return provider, nil
}

// ProviderConstructor is useful when client wants to add their own custom
// TracerProvider.
type ProviderConstructor func(config Config) (trace.TracerProvider, error)

// Created pre-defined immutable map of built-in provider's
var providersConfig = map[string]ProviderConstructor{
	"jaeger": func(cfg Config) (trace.TracerProvider, error) {
		traceProvider, _, err := jaeger.NewExportPipeline(
			jaeger.WithCollectorEndpoint(cfg.Endpoint),
			jaeger.WithSDKOptions(
				sdktrace.WithSampler(sdktrace.AlwaysSample()),
				sdktrace.WithResource(
					resource.NewWithAttributes(
						semconv.ServiceNameKey.String(cfg.ApplicationName),
						attribute.String("exporter", cfg.Provider),
					),
				),
			),
		)
		if err != nil {
			return nil, err
		}
		return traceProvider, nil
	},
	"zipkin": func(cfg Config) (trace.TracerProvider, error) {
		traceProvider, err := zipkin.NewExportPipeline(cfg.Endpoint,
			zipkin.WithSDKOptions(
				sdktrace.WithSampler(sdktrace.AlwaysSample()),
				sdktrace.WithResource(
					resource.NewWithAttributes(semconv.ServiceNameKey.String(cfg.ApplicationName)),
				),
			),
		)
		return traceProvider, err
	},
	"stdout": func(cfg Config) (trace.TracerProvider, error) {
		var option stdout.Option
		if cfg.SkipTraceExport {
			option = stdout.WithoutTraceExport()
		} else {
			option = stdout.WithPrettyPrint()
		}
		otExporter, err := stdout.NewExporter(option)
		if err != nil {
			return nil, err
		}
		traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(otExporter))
		return traceProvider, nil
	},
	"noop": func(config Config) (trace.TracerProvider, error) {
		return trace.NewNoopTracerProvider(), nil
	},
}
