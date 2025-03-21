package usecases

import (
	"context"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"labs-two-service-a/internal/infra/services"
	"labs-two-service-a/internal/infra/tracing"
)

type GetCepUseCaseInterface interface {
	GetTempoPorCep(ctx context.Context, cep string) (entities.TempoResponseDto, error)
}

type GetCepUseCase struct {
	appConfid              *config.AppSettings
	ServiceTempo           services.ServiceTempoInterface
	ViaCepServiceInterface services.ServiceCepInterface
	tracingProvider        *tracing.TracingProvider
}

func NewGetCepUseCase(appConfig *config.AppSettings, viaCepService services.ServiceCepInterface, tempoService services.ServiceTempoInterface, tracingProvider *tracing.TracingProvider) *GetCepUseCase {
	return &GetCepUseCase{
		appConfid:              appConfig,
		ServiceTempo:           tempoService,
		ViaCepServiceInterface: viaCepService,
		tracingProvider:        tracingProvider,
	}
}

func (u *GetCepUseCase) GetTempoPorCep(ctx context.Context, cep string) (entities.TempoResponseDto, error) {
	isValidCep := ValidateCEP(cep)
	if !isValidCep {
		return entities.TempoResponseDto{}, &entities.CustomErrors{
			Code:    422,
			Message: "invalid zipcode",
		}
	}

	ctxCep, spanCep := u.tracingProvider.Tracer.Start(ctx, "GetCep")
	defer spanCep.End()

	cepResponse, err := u.ViaCepServiceInterface.GetCep(ctxCep, cep)
	if err != nil {
		spanCep.RecordError(err)
		return entities.TempoResponseDto{}, &entities.CustomErrors{
			Code:    404,
			Message: "can not find zipcode",
		}
	}

	ctxTempo, spanTempo := u.tracingProvider.Tracer.Start(ctx, "GetTempo")
	defer spanTempo.End()
	
	tempoResponse, err := u.ServiceTempo.GetTempo(ctxTempo, cepResponse.Localidade)
	if err != nil {
		spanTempo.RecordError(err)
		return entities.TempoResponseDto{}, &entities.CustomErrors{
			Code:    404,
			Message: "can not find weather",
		}
	}

	return tempoResponse, nil
}

func ValidateCEP(cep string) bool {
	isValid := len(cep) == 8
	return isValid
}
