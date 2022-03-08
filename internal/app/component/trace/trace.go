package trace

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
)

type Config struct {
	Endpoint    string
	ServiceName string
	Env         string
	Timeout     int64
}

type Tracer struct {
	config         Config
	tracerProvider *sdktrace.TracerProvider
}

// New 创建 Tracer
func New(config *Config, logger log.Logger) (*Tracer, func(), error) {
	if config == nil {
		return nil, func() {}, nil
	}

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
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.DeploymentEnvironmentKey.String(config.Env),
		)),
	)

	cleanup := func() {
		log.NewHelper(logger).Info("closing the trace")

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.NewHelper(logger).Error(err)
		}
	}

	return &Tracer{
		config:         *config,
		tracerProvider: tp,
	}, cleanup, nil
}

func (t *Tracer) TracerProvider() *sdktrace.TracerProvider {
	return t.tracerProvider
}

func (t *Tracer) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	if name == "" {
		name = t.config.ServiceName
	}
	return t.tracerProvider.Tracer(name, opts...)
}
