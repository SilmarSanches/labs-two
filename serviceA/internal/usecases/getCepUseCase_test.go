package usecases

import (
	"context"
	"errors"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockViaCepService struct {
	mock.Mock
}

func (m *MockViaCepService) GetCep(ctx context.Context, cep string) (entities.ViaCepDto, error) {
	args := m.Called(ctx, cep)
	return args.Get(0).(entities.ViaCepDto), args.Error(1)
}

func TestGetCep_Success(t *testing.T) {
	// Arrange
	mockViaCepService := new(MockViaCepService)
	appConfig := &config.AppSettings{}

	useCase := NewGetCepUseCase(appConfig, mockViaCepService)

	cep := "01001-000"
	cidade := "São Paulo"
	mockViaCepService.On("GetCep", mock.Anything, cep).Return(entities.ViaCepDto{Localidade: cidade}, nil)

	// Act
	result, err := useCase.GetCep(cep)

	// Assert
	require.NoError(t, err)
	require.Equal(t, cidade, result.Localidade)
	mockViaCepService.AssertCalled(t, "GetCep", mock.Anything, cep)
}

func TestGetTempo_CepNotFount(t *testing.T) {
	// Arrange
	mockViaCepService := new(MockViaCepService)
	appConfig := &config.AppSettings{}

	useCase := NewGetCepUseCase(appConfig, mockViaCepService)

	cep := "01001-000"
	mockViaCepService.On("GetCep", mock.Anything, cep).Return(entities.ViaCepDto{}, errors.New("CEP not found"))

	// Act
	_, err := useCase.GetCep(cep)

	// Assert
	require.Error(t, err)
	require.IsType(t, &entities.CustomErrors{}, err)
	customErr := err.(*entities.CustomErrors)
	require.Equal(t, 404, customErr.Code)
	require.Equal(t, "can not find zipcode", customErr.Message)
	mockViaCepService.AssertCalled(t, "GetCep", mock.Anything, cep)
}