package tracer

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestConfigureTracerProvider(t *testing.T) {
	tcs := []struct {
		Description string
		Config      Config
		Err         error
	}{
		{
			Description: "Jaeger: Missing endpoint",
			Config: Config{
				Provider: "jaeger",
			},
			Err: ErrTracerProviderBuildFailed,
		},
		{
			Description: "Zipkin: Missing endpoint",
			Config: Config{
				Provider: "Zipkin",
			},
			Err: ErrTracerProviderBuildFailed,
		},
		{
			Description: "Jaeger: Valid",
			Config: Config{
				Provider: "jaeger",
				Endpoint: "http://localhost",
			},
		},
		{
			Description: "Zipkin: Valid",
			Config: Config{
				Provider: "Zipkin",
				Endpoint: "http://localhost",
			},
		},
		{
			Description: "Unknown Provider",
			Config: Config{
				Provider: "undefined",
			},
			Err: ErrTracerProviderNotFound,
		},
		{
			Description: "Stdout: Valid",
			Config: Config{
				Provider: "stdOut",
			},
		},
		{
			Description: "Stdout: Valid skip export",
			Config: Config{
				Provider:        "stdoUt",
				SkipTraceExport: true,
			},
		},
		{
			Description: "Default",
			Config:      Config{},
		},
		{
			Description: "NoOp: Valid",
			Config: Config{
				Provider: "noop",
			},
		},
		{
			Description: "Custom provider",
			Config: Config{
				Provider: "coolest",
				Providers: map[string]ProviderConstructor{
					"coolest": func(_ Config) (trace.TracerProvider, error) {
						return trace.NewNoopTracerProvider(), nil
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				assert  = assert.New(t)
				tp, err = ConfigureTracerProvider(tc.Config)
			)
			if tc.Err == nil {
				assert.NotNil(tp)
			}
			assert.True(errors.Is(err, tc.Err))
		})
	}
}
