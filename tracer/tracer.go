package tracer

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

// ExtractTraceInfo returns the ID of the trace flowing through the context
// as well as ID the current active span. The third boolean return value represents
// whether the returned IDs are valid and safe to use. OpenTelemetry's noop
// tracer provider, for instance, returns zero value trace information that's
// considered invalid and should be ignored.
func ExtractTraceInfo(ctx context.Context) (string, string, bool) {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	return traceID, spanID, span.SpanContext().IsValid()
}

