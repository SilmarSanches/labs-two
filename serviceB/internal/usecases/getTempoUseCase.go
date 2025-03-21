package usecases

import (
	"context"
	"labs-two-service-b/config"
	"labs-two-service-b/internal/entities"
	"labs-two-service-b/internal/infra/services"
	"labs-two-service-b/internal/infra/tracing"
)

type GetTempoUseCaseInterface interface {
	GetTempo(ctx context.Context, cep string) (entities.GetTempoResponseDto, error)
}

type GetTempoUseCase struct {
	appConfid               *config.AppSettings
	tracingProvider         *tracing.TracingProvider
	WeatherServiceInterface services.ServiceTempoInterface
}

func NewGetTempoUseCase(appConfig *config.AppSettings, weatherService services.ServiceTempoInterface, tracingProvider *tracing.TracingProvider) *GetTempoUseCase {
	return &GetTempoUseCase{
		appConfid:               appConfig,
		WeatherServiceInterface: weatherService,
		tracingProvider:        tracingProvider,
	}
}

func (u *GetTempoUseCase) GetTempo(ctx context.Context, location string) (entities.GetTempoResponseDto, error) {
	isValidLocation := ValidateLocation(location)
	if !isValidLocation {
		return entities.GetTempoResponseDto{}, &entities.CustomErrors{
			Code:    422,
			Message: "invalid location",
		}
	}

	ctxWeather, spanWeather := u.tracingProvider.Tracer.Start(ctx, "GetWeather")
	defer spanWeather.End()

	weather, err := u.WeatherServiceInterface.GetTempo(ctxWeather, location)
	if err != nil {
		spanWeather.RecordError(err)
		return entities.GetTempoResponseDto{}, &entities.CustomErrors{
			Code:    404,
			Message: "can not find location",
		}
	}

	celcius := weather.Current.TempC
	Kelvin := celcius + 273
	Fahrenheit := celcius*1.8 + 32

	result := entities.GetTempoResponseDto{
		Kelvin:     Kelvin,
		Celsius:    celcius,
		Fahrenheit: Fahrenheit,
		City:       location,
	}

	return result, nil
}

func ValidateLocation(location string) bool {
	return len(location) >= 4
}
