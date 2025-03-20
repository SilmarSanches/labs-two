package tracing

import (
	"context"
	"labs-two-service-b/config"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var tracingInstance *TracingProvider

type TracingConfig struct {
	ServiceName string
	ZipkinURL   string
}

type TracingProvider struct {
	Tracer trace.Tracer
	tp     *sdktrace.TracerProvider
}

func ProvideTracingConfig(cfg *config.AppSettings) TracingConfig {
	return TracingConfig{
		ServiceName: cfg.ServiceName,
		ZipkinURL:   cfg.UrlZipkin,
	}
}

func ProvideTracingProvider(cfg TracingConfig) *TracingProvider {
	if tracingInstance != nil {
		return tracingInstance
	}

	provider, _, err := NewTracingProvider(cfg)
	if err != nil {
		panic(err)
	}

	tracingInstance = provider
	return provider
}

func ProvideTracingProviderWithCleanup(cfg TracingConfig) (*TracingProvider, func()) {
	if tracingInstance != nil {
		return tracingInstance, func() {}
	}

	provider, cleanup, err := NewTracingProvider(cfg)
	if err != nil {
		panic(err)
	}

	tracingInstance = provider
	return provider, cleanup
}

func NewTracingProvider(cfg TracingConfig) (*TracingProvider, func(), error) {
	exporter, err := zipkin.New(cfg.ZipkinURL)
	if err != nil {
		log.Fatalf("Erro ao criar o Exporter Zipkin: %v", err)
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)),
	)

	otel.SetTracerProvider(tp)

	tracingProvider := &TracingProvider{
		Tracer: tp.Tracer(cfg.ServiceName),
		tp:     tp,
	}

	cleanup := func() {
		_ = tp.Shutdown(context.Background())
	}

	return tracingProvider, cleanup, nil
}