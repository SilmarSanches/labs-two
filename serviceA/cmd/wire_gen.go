// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/infra/services"
	"labs-two-service-a/internal/infra/tracing"
	"labs-two-service-a/internal/infra/web"
	"labs-two-service-a/internal/usecases"
)

// Injectors from wire.go:

func NewConfig() *config.AppSettings {
	appSettings := config.ProvideConfig()
	return appSettings
}

func NewGetConsultaUseCase() *usecases.GetConsultaUseCase {
	appSettings := config.ProvideConfig()
	httpClient := services.NewHttpClient()
	serviceConsulta := services.NewServiceConsulta(httpClient, appSettings)
	tracingConfig := tracing.ProvideTracingConfig(appSettings)
	tracingProvider := tracing.ProvideTracingProvider(tracingConfig)
	getConsultaUseCase := usecases.NewGetConsultaUseCase(appSettings, serviceConsulta, tracingProvider)
	return getConsultaUseCase
}

func NewGetConsultaHandler() *web.GetConsultaHandler {
	appSettings := config.ProvideConfig()
	httpClient := services.NewHttpClient()
	serviceConsulta := services.NewServiceConsulta(httpClient, appSettings)
	tracingConfig := tracing.ProvideTracingConfig(appSettings)
	tracingProvider := tracing.ProvideTracingProvider(tracingConfig)
	getConsultaUseCase := usecases.NewGetConsultaUseCase(appSettings, serviceConsulta, tracingProvider)
	getConsultaHandler := web.NewGetConsultaHandler(appSettings, getConsultaUseCase, serviceConsulta, tracingProvider)
	return getConsultaHandler
}

func InitializeTracing() (*tracing.TracingProvider, func()) {
	appSettings := config.ProvideConfig()
	tracingConfig := tracing.ProvideTracingConfig(appSettings)
	tracingProvider, cleanup := tracing.ProvideTracingProviderWithCleanup(tracingConfig)
	return tracingProvider, func() {
		cleanup()
	}
}

// wire.go:

var ProviderConfig = wire.NewSet(config.ProvideConfig)

var ProviderHttpClient = wire.NewSet(services.NewHttpClient)

var ProviderConsulta = wire.NewSet(services.NewServiceConsulta, wire.Bind(new(services.ServiceConsultaInterface), new(*services.ServiceConsulta)))

var ProviderTracingForHandler = wire.NewSet(tracing.ProvideTracingConfig, tracing.ProvideTracingProvider)

var ProviderTracingWithCleanup = wire.NewSet(tracing.ProvideTracingConfig, tracing.ProvideTracingProviderWithCleanup)

var ProviderGlobal = wire.NewSet(
	ProviderHttpClient,
	ProviderConfig,
	ProviderConsulta,
	ProviderTracingForHandler,
)

var ProviderUseCase = wire.NewSet(usecases.NewGetConsultaUseCase, wire.Bind(new(usecases.GetConsultaUseCaseInterface), new(*usecases.GetConsultaUseCase)))

var ProviderHandler = wire.NewSet(web.NewGetConsultaHandler)
