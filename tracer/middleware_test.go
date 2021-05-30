package tracer

import (
	"net/http"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/propagation"
)

func TestEchoFirstTraceNodeInfo(t *testing.T) {
	type args struct {
		propagator propagation.TextMapPropagator
	}
	tests := []struct {
		name string
		args args
		want func(http.Handler) http.Handler
	}{
		{"TestEchoFirstTraceNodeInfo",
		nil,
		nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EchoFirstTraceNodeInfo(tt.args.propagator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EchoFirstTraceNodeInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
