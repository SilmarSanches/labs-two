//go:build wireinject
// +build wireinject

package main

import (
	"labs-two-service-a/config"
	"labs-two-service-a/internal/infra/services"
	"labs-two-service-a/internal/infra/tracing"
	"labs-two-service-a/internal/infra/web"
	"labs-two-service-a/internal/usecases"

	"github.com/google/wire"
)

var ProviderConfig = wire.NewSet(config.ProvideConfig)

var ProviderHttpClient = wire.NewSet(
	services.NewHttpClient,
)

var ProviderConsulta = wire.NewSet(
	services.NewServiceConsulta,
	wire.Bind(new(services.ServiceConsultaInterface), new(*services.ServiceConsulta)),
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
	ProviderConsulta,
	ProviderTracingForHandler,
)

var ProviderUseCase = wire.NewSet(
	usecases.NewGetConsultaUseCase,
	wire.Bind(new(usecases.GetConsultaUseCaseInterface), new(*usecases.GetConsultaUseCase)),
)

var ProviderHandler = wire.NewSet(
	web.NewGetConsultaHandler,
)

func NewConfig() *config.AppSettings {
	wire.Build(ProviderConfig)
	return &config.AppSettings{}
}

func NewGetConsultaUseCase() *usecases.GetConsultaUseCase {
	wire.Build(ProviderGlobal, ProviderUseCase)
	return &usecases.GetConsultaUseCase{}
}

func NewGetConsultaHandler() *web.GetConsultaHandler {
	wire.Build(ProviderGlobal, ProviderUseCase, ProviderHandler)
	return &web.GetConsultaHandler{}
}

func InitializeTracing() (*tracing.TracingProvider, func()) {
	wire.Build(ProviderConfig, ProviderTracingWithCleanup)
	return nil, nil
}
