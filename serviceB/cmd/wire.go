//go:build wireinject
// +build wireinject

package main

import (
	"labs-two-service-b/config"
	"labs-two-service-b/internal/infra/services"
	"labs-two-service-b/internal/infra/web"
	"labs-two-service-b/internal/usecases"

	"github.com/google/wire"
)

var ProviderConfig = wire.NewSet(config.ProvideConfig)

var ProviderHttpClient = wire.NewSet(
	services.NewHttpClient,
)

var ProviderTempo = wire.NewSet(
	services.NewServiceTempo,
	wire.Bind(new(services.ServiceTempoInterface), new(*services.ServiceTempo)),
)

var ProviderGlobal = wire.NewSet(
	ProviderHttpClient,
	ProviderConfig,
	ProviderTempo,
)

var ProviderUseCase = wire.NewSet(
	usecases.NewGetTempoUseCase,
	wire.Bind(new(usecases.GetTempoUseCaseInterface), new(*usecases.GetTempoUseCase)),
)

var ProviderHandler = wire.NewSet(web.NewGetCepHandler)

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
