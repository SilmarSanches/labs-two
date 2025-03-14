package usecases

import (
	"context"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"labs-two-service-a/internal/infra/services"
	"regexp"
)

type GetCepUseCaseInterface interface {
	GetCep(cep string) (entities.ViaCepDto, error)
}

type GetCepUseCase struct {
	appConfid              *config.AppSettings
	ViaCepServiceInterface services.ServiceCepInterface
}

func NewGetCepUseCase(appConfig *config.AppSettings, viaCepService services.ServiceCepInterface) *GetCepUseCase {
	return &GetCepUseCase{
		appConfid:              appConfig,
		ViaCepServiceInterface: viaCepService,
	}
}

func (u *GetCepUseCase) GetCep(cep string) (entities.ViaCepDto, error) {
	ctx := context.Background()

	isValidCep := ValidateCEP(cep)
	if !isValidCep {
		return entities.ViaCepDto{}, &entities.CustomErrors{
			Code:    422,
			Message: "invalid zipcode",
		}
	}

	cepResponse, err := u.ViaCepServiceInterface.GetCep(ctx, cep)
	if err != nil {
		return entities.ViaCepDto{}, &entities.CustomErrors{
			Code:    404,
			Message: "can not find zipcode",
		}
	}

	return cepResponse, nil
}

func ValidateCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{5}-\d{3}$`)
	return re.MatchString(cep)
}
