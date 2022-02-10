package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"strings"
)

// New 返回 *trace.TracerProvider
func New(config Config) (*trace.TracerProvider, error) {
	var endpointOption jaeger.EndpointOption

	if strings.HasPrefix(config.Endpoint, "http") {
		endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Endpoint))
	} else {
		agentConfig := strings.SplitN(config.Endpoint, ":", 2)
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

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.DeploymentEnvironmentKey.String(config.Env),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

// MustNew 返回 *trace.TracerProvider
func MustNew(config Config) *trace.TracerProvider {
	tp, err := New(config)
	if err != nil {
		panic(err)
	}

	return tp
}
