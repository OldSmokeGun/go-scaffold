package trace

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	tlog "go-scaffold/pkg/log/otel/trace"
)

var ErrUnsupportedOTPLProtocol = errors.New("unsupported otpl protocol")

// OTPLProtocol OTPL protocol
type OTPLProtocol string

const (
	HTTP OTPLProtocol = "http"
	GRPC OTPLProtocol = "grpc"
)

// Trace OpenTelemetry trace
type Trace struct {
	serviceName    string
	env            string
	tracerProvider *sdktrace.TracerProvider
	errorLogger    *tlog.ErrorLogger
}

type Option func(t *Trace)

// WithServiceName optional service name
func WithServiceName(serviceName string) Option {
	return func(t *Trace) {
		t.serviceName = serviceName
	}
}

// WithEnv optional environment value
func WithEnv(env string) Option {
	return func(t *Trace) {
		t.env = env
	}
}

// WithErrorLogger optional error logger
func WithErrorLogger(logger *tlog.ErrorLogger) Option {
	return func(t *Trace) {
		t.errorLogger = logger
	}
}

// New build Trace
func New(ctx context.Context, protocol OTPLProtocol, endpoint string, options ...Option) (*Trace, error) {
	t := &Trace{}
	for _, option := range options {
		option(t)
	}

	if t.errorLogger != nil {
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			t.errorLogger.Handle(err)
		}))
	}

	var (
		exporter *otlptrace.Exporter
		err      error
	)

	switch protocol {
	case HTTP:
		exporter, err = otlptracehttp.New(
			ctx,
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(endpoint),
		)
		if err != nil {
			return nil, err
		}
	case GRPC:
		exporter, err = otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(endpoint),
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedOTPLProtocol
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(t.serviceName),
			semconv.DeploymentEnvironmentKey.String(t.env),
		)),
	)

	t.tracerProvider = tp

	return t, nil
}

// TracerProvider return OpenTelemetry TracerProvider
func (t *Trace) TracerProvider() *sdktrace.TracerProvider {
	return t.tracerProvider
}

// Tracer return a OpenTelemetry Tracer
func (t *Trace) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	if name == "" {
		name = t.serviceName
	}
	return t.tracerProvider.Tracer(name, opts...)
}

// Shutdown the span processors
func (t *Trace) Shutdown(ctx context.Context) error {
	if t.tracerProvider == nil {
		return nil
	}
	return t.tracerProvider.Shutdown(ctx)
}
