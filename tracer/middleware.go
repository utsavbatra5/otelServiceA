package tracer

import (
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	spanIDHeaderName  = "W3C-Span-ID"
	traceIDHeaderName = "W3C-Trace-ID"
	SpanIDLogKeyName  = "span-id"
	TraceIdLogKeyName = "trace-id"
)

// EchoFirstNodeTraceInfo captures the trace information from a request and writes it
// back in the response headers if the request is the first one in the trace path.
func EchoFirstTraceNodeInfo(propagator propagation.TextMapPropagator) func(http.Handler) http.Handler {
	return func(delegate http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			rsc := trace.RemoteSpanContextFromContext(ctx)
			sc := trace.SpanContextFromContext(ctx)
			if sc.IsValid() && !rsc.IsValid() {
				w.Header().Set(spanIDHeaderName, sc.SpanID().String())
				w.Header().Set(traceIDHeaderName, sc.TraceID().String())
			}
			delegate.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
