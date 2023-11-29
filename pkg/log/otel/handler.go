package otel

import (
	"context"
	"log/slog"

	sdktrace "go.opentelemetry.io/otel/trace"
)

const (
	TraceIDLogKey = "trace_id"
	SpanIDLogKey  = "span_id"
)

// Handler slog handler
type Handler struct {
	slog.Handler
}

func NewHandler(handler slog.Handler) *Handler {
	return &Handler{Handler: handler}
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	spanCtx := sdktrace.SpanContextFromContext(ctx)

	var attrs []slog.Attr

	if spanCtx.HasTraceID() {
		attrs = append(attrs, slog.String(TraceIDLogKey, spanCtx.TraceID().String()))
	}

	if spanCtx.HasSpanID() {
		attrs = append(attrs, slog.String(SpanIDLogKey, spanCtx.SpanID().String()))
	}

	if len(attrs) > 0 {
		record.Add("trace", slog.GroupValue(attrs...))
	}

	return h.Handler.Handle(ctx, record)
}
