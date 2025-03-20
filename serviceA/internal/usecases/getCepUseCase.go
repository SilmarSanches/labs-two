package usecases

import (
	"context"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"labs-two-service-a/internal/infra/services"
	"regexp"
)

type GetCepUseCaseInterface interface {
	GetTempoPorCep(cep string) (entities.TempoResponseDto, error)
}

type GetCepUseCase struct {
	appConfid              *config.AppSettings
	ServiceTempo           services.ServiceTempoInterface
	ViaCepServiceInterface services.ServiceCepInterface
}

func NewGetCepUseCase(appConfig *config.AppSettings, viaCepService services.ServiceCepInterface, tempoService services.ServiceTempoInterface) *GetCepUseCase {
	return &GetCepUseCase{
		appConfid:              appConfig,
		ServiceTempo:           tempoService,
		ViaCepServiceInterface: viaCepService,
	}
}

func (u *GetCepUseCase) GetTempoPorCep(cep string) (entities.TempoResponseDto, error) {
	ctx := context.Background()

	isValidCep := ValidateCEP(cep)
	if !isValidCep {
		return entities.TempoResponseDto{}, &entities.CustomErrors{
			Code:    422,
			Message: "invalid zipcode",
		}
	}

	cepResponse, err := u.ViaCepServiceInterface.GetCep(ctx, cep)
	if err != nil {
		return entities.TempoResponseDto{}, &entities.CustomErrors{
			Code:    404,
			Message: "can not find zipcode",
		}
	}

	tempoResponse, err := u.ServiceTempo.GetTempo(ctx, cepResponse.Localidade)
	if err != nil {
		return entities.TempoResponseDto{}, &entities.CustomErrors{
			Code:    404,
			Message: "can not find weather",
		}
	}

	return tempoResponse, nil
}

func ValidateCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{5}-\d{3}$`)
	return re.MatchString(cep)
}
