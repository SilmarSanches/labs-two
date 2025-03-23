package usecases

import (
	"context"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"labs-two-service-a/internal/infra/services"
	"labs-two-service-a/internal/infra/tracing"
)

type GetConsultaUseCaseInterface interface {
	GetTempoPorCep(ctx context.Context, cep entities.CepRequestDto) (entities.TempoResponseDto, error)
}

type GetConsultaUseCase struct {
	appConfid       *config.AppSettings
	ServiceConsulta services.ServiceConsultaInterface
	tracingProvider *tracing.TracingProvider
}

func NewGetConsultaUseCase(appConfig *config.AppSettings, serviceConsulta services.ServiceConsultaInterface, tracingProvider *tracing.TracingProvider) *GetConsultaUseCase {
	return &GetConsultaUseCase{
		appConfid:       appConfig,
		ServiceConsulta: serviceConsulta,
		tracingProvider: tracingProvider,
	}
}

func (u *GetConsultaUseCase) GetTempoPorCep(ctx context.Context, cep entities.CepRequestDto) (entities.TempoResponseDto, error) {
	isValidCep := ValidateCEP(cep.Cep)
	if !isValidCep {
		return entities.TempoResponseDto{}, &entities.CustomErrors{
			Code:    422,
			Message: "invalid zipcode",
		}
	}

	ctxTempo, spanTempo := u.tracingProvider.Tracer.Start(ctx, "GetCep")
	defer spanTempo.End()

	tempoResponse, err := u.ServiceConsulta.GetTempo(ctxTempo, cep)
	if err != nil {
		spanTempo.RecordError(err)
		return entities.TempoResponseDto{}, &entities.CustomErrors{
			Code:    404,
			Message: "can not find zipcode",
		}
	}

	return tempoResponse, nil
}

func ValidateCEP(cep string) bool {
	isValid := len(cep) == 8
	return isValid
}
