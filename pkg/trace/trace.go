package trace

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// Trace OpenTelemetry trace
type Trace struct {
	endpoint       string
	serviceName    string
	env            string
	tracerProvider *sdktrace.TracerProvider
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

// New build Trace
func New(endpoint string, options ...Option) (*Trace, error) {
	t := &Trace{}
	for _, option := range options {
		option(t)
	}

	var endpointOption jaeger.EndpointOption

	if strings.HasPrefix(endpoint, "http") {
		endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
	} else {
		agentConfig := strings.SplitN(endpoint, ":", 2)
		if len(agentConfig) == 2 {
			endpointOption = jaeger.WithAgentEndpoint(
				jaeger.WithAgentHost(agentConfig[0]),
				jaeger.WithAgentPort(agentConfig[1]),
			)
		} else {
			endpointOption = jaeger.WithAgentEndpoint(jaeger.WithAgentHost(agentConfig[0]))
		}
	}

	exporter, err := jaeger.New(endpointOption)
	if err != nil {
		return nil, err
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
