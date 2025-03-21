//go:build wireinject
// +build wireinject

package main

import (
	"labs-two-serviceb/config"
	"labs-two-serviceb/internal/infra/services"
    "labs-two-serviceb/internal/infra/tracing"
	"labs-two-serviceb/internal/infra/web"
	"labs-two-serviceb/internal/usecases"

	"github.com/google/wire"
)

var ProviderConfig = wire.NewSet(config.ProvideConfig)

var ProviderHttpClient = wire.NewSet(
    services.NewHttpClient,
)

var ProviderCep = wire.NewSet(
    services.NewServiceCep,
    wire.Bind(new(services.ServiceCepInterface), new(*services.ServiceCep)),
)

var ProviderTempo = wire.NewSet(
    services.NewServiceTempo,
    wire.Bind(new(services.ServiceTempoInterface), new(*services.ServiceTempo)),
)

var ProviderTracingForHandler = wire.NewSet(
	tracing.ProvideTracingConfig,
	tracing.ProvideTracingProvider,
)

var ProviderTracingWithCleanup = wire.NewSet(
	tracing.ProvideTracingConfig,
	tracing.ProvideTracingProviderWithCleanup,
)

var ProviderGlobal = wire.NewSet(
    ProviderHttpClient,
    ProviderConfig,
    ProviderCep,
    ProviderTempo,
    ProviderTracingForHandler,
)

var ProviderUseCase = wire.NewSet(
	usecases.NewGetTempoUseCase,
	wire.Bind(new(usecases.GetTempoUseCaseInterface), new(*usecases.GetTempoUseCase)),
)

var ProviderHandler = wire.NewSet(web.NewGetTempoHandler)

func NewConfig() *config.AppSettings {
	wire.Build(ProviderConfig)
	return &config.AppSettings{}
}

func NewGetTempUseCase() *usecases.GetTempoUseCase {
    wire.Build(ProviderGlobal, ProviderUseCase)
    return &usecases.GetTempoUseCase{}
}

func NewGetTempoHandler() *web.GetTempoHandler {
    wire.Build(ProviderGlobal, ProviderUseCase, ProviderHandler)
    return &web.GetTempoHandler{}
}

func InitializeTracing() (*tracing.TracingProvider, func()) {
	wire.Build(ProviderConfig, ProviderTracingWithCleanup)
	return nil, nil
}