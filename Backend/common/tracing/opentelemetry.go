package tracing

import (
	"context"
	"io"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitTracer(serviceName, agentHost, agentPort string) (io.Closer, error) {
	log.Printf("Initializing tracing to Jaeger Agent at %s:%s\n", agentHost, agentPort)

	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(agentHost),
		jaeger.WithAgentPort(agentPort),
	))
	if err != nil {
		log.Fatalf("failed to create Jaeger exporter: %v", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &tracerProviderCloser{tp}, nil
}

type tracerProviderCloser struct {
	provider *sdktrace.TracerProvider
}

func (c *tracerProviderCloser) Close() error {
	if err := c.provider.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
		return err
	}
	return nil
}
