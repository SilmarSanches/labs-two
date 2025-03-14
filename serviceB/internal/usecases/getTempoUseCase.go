package usecases

import (
	"context"
	"labs-two-service-b/config"
	"labs-two-service-b/internal/entities"
	"labs-two-service-b/internal/infra/services"
)

type GetTempoUseCaseInterface interface {
	GetTempo(cep string) (entities.GetTempoResponseDto, error)
}

type GetTempoUseCase struct {
	appConfid               *config.AppSettings
	WeatherServiceInterface services.ServiceTempoInterface
}

func NewGetTempoUseCase(appConfig *config.AppSettings, weatherService services.ServiceTempoInterface) *GetTempoUseCase {
	return &GetTempoUseCase{
		appConfid:               appConfig,
		WeatherServiceInterface: weatherService,
	}
}

func (u *GetTempoUseCase) GetTempo(location string) (entities.GetTempoResponseDto, error) {
	ctx := context.Background()

	isValidLocation := ValidateLocation(location)
	if !isValidLocation {
		return entities.GetTempoResponseDto{}, &entities.CustomErrors{
			Code:    422,
			Message: "invalid location",
		}
	}

	weather, err := u.WeatherServiceInterface.GetTempo(ctx, location)
	if err != nil {
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
