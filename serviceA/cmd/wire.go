// go:build wireinject
// +build wireinject

package main

import (
	"labs-two-service-a/config"
	"labs-two-service-a/internal/infra/services"
	"labs-two-service-a/internal/infra/web"
	"labs-two-service-a/internal/usecases"

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

var ProviderGlobal = wire.NewSet(
    ProviderHttpClient,
    ProviderConfig,
    ProviderTempo,
    ProviderCep,
)

var ProviderUseCase = wire.NewSet(
	usecases.NewGetCepUseCase,
	wire.Bind(new(usecases.GetCepUseCaseInterface), new(*usecases.GetCepUseCase)),
)

var ProviderHandler = wire.NewSet(web.NewGetCepHandler)

func NewConfig() *config.AppSettings {
	wire.Build(ProviderConfig)
	return &config.AppSettings{}
}

func NewGetCepUseCase() *usecases.GetCepUseCase {
    wire.Build(ProviderGlobal, ProviderUseCase)
    return &usecases.GetCepUseCase{}
}

func NewGetCepHandler() *web.GetCepHandler {
    wire.Build(ProviderGlobal, ProviderUseCase, ProviderHandler)
    return &web.GetCepHandler{}
}