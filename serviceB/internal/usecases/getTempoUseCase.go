package usecases

import (
	"context"
	"labs-two-serviceb/config"
	"labs-two-serviceb/internal/entities"
	"labs-two-serviceb/internal/infra/services"
	"labs-two-serviceb/internal/infra/tracing"
	"regexp"
)

type GetTempoUseCaseInterface interface {
	GetTempo(ctx context.Context, cep string) (entities.GetTempoResponseDto, error)
}

type GetTempoUseCase struct {
	appConfid               *config.AppSettings
	tracingProvider         *tracing.TracingProvider
	ViaCepServiceInterface  services.ServiceCepInterface
	WeatherServiceInterface services.ServiceTempoInterface
}

func NewGetTempoUseCase(appConfig *config.AppSettings, viaCepService services.ServiceCepInterface, weatherService services.ServiceTempoInterface, tracingProvider *tracing.TracingProvider) *GetTempoUseCase {
	return &GetTempoUseCase{
		appConfid:               appConfig,
		ViaCepServiceInterface:  viaCepService,
		WeatherServiceInterface: weatherService,
		tracingProvider:        tracingProvider,
	}
}

func (u *GetTempoUseCase) GetTempo(ctx context.Context, cep string) (entities.GetTempoResponseDto, error) {

	isValidCep := len(cep) == 8 && regexp.MustCompile(`^[0-9]+$`).MatchString(cep)
	if !isValidCep {
		return entities.GetTempoResponseDto{}, &entities.CustomError{
			Code:    422,
			Message: "invalid zipcode",
		}
	}

	ctxCep, spanCep := u.tracingProvider.Tracer.Start(ctx, "GetCep")
	defer spanCep.End()

	cepResponse, err := u.ViaCepServiceInterface.GetCep(ctxCep, cep)
	if err != nil {
		spanCep.RecordError(err)
		return entities.GetTempoResponseDto{}, &entities.CustomError{
			Code:    404,
			Message: "can not find zipcode",
		}
	}

	ctxTempo, spanTempo := u.tracingProvider.Tracer.Start(ctx, "GetTempo")
	defer spanTempo.End()

	weather, err := u.WeatherServiceInterface.GetTempo(ctxTempo, cepResponse.Localidade)
	if err != nil {
		spanTempo.RecordError(err)
		return entities.GetTempoResponseDto{}, &entities.CustomError{
			Code:    404,
			Message: "can not find temperature",
		}
	}

	celcius := weather.Current.TempC
	Kelvin := celcius + 273
	Fahrenheit := celcius*1.8 + 32

	result := entities.GetTempoResponseDto{
		Kelvin:     Kelvin,
		Celsius:    celcius,
		Fahrenheit: Fahrenheit,
		City:       cepResponse.Localidade,
	}

	return result, nil
}
