package tracing

import (
	"context"
	"io"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
)

func InitTracer(serviceName, otlpEndpoint string) (io.Closer, error) {
	log.Printf("Initializing tracing to OTLP endpoint at %s\n", otlpEndpoint)

	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		return nil, err
	}

	// === ISPRAVKA JE OVDE ===
	// Kreiramo resurs sa odvojenim opcijama za SchemaURL i za atribute.
	res, err := resource.New(context.Background(),
		resource.WithSchemaURL(semconv.SchemaURL), // SchemaURL se postavlja ovako
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName), // A ostali atributi idu ovde
		),
	)
	if err != nil {
		return nil, err
	}
	// === KRAJ ISPRAVKE ===

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
