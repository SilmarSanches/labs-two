package usecases

import (
	"context"
	"errors"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"labs-two-service-a/internal/infra/tracing"
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

type MockTempoService struct {
	mock.Mock
}

func (m *MockTempoService) GetTempo(ctx context.Context, cidade string) (entities.TempoResponseDto, error) {
	args := m.Called(ctx, cidade)
	return args.Get(0).(entities.TempoResponseDto), args.Error(1)
}

func TestGetCep_Success(t *testing.T) {
	// Arrange
	mockViaCepService := new(MockViaCepService)
	mockTempoService := new(MockTempoService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetCepUseCase(appConfig, mockViaCepService, mockTempoService, tracingProvider)

	cep := "01001-000"
	cidade := "São Paulo"
	mockViaCepService.On("GetCep", mock.Anything, cep).Return(entities.ViaCepDto{Localidade: cidade}, nil)
	mockTempoService.On("GetTempo", mock.Anything, cidade).Return(entities.TempoResponseDto{
		City: cidade,
		Kelvin: 273.15,
		Celsius: 0,
		Fahrenheit: 32,
	}, nil)

	// Act
	result, err := useCase.GetTempoPorCep(context.Background(),cep)

	// Assert
	require.NoError(t, err)
	require.Equal(t, cidade, result.City)
	require.Equal(t, 273.15, result.Kelvin)
	require.Equal(t, 0.0, result.Celsius)
	require.Equal(t, 32.0, result.Fahrenheit)
	mockViaCepService.AssertCalled(t, "GetCep", mock.Anything, cep)
}

func TestGetTempo_CepNotFount(t *testing.T) {
	// Arrange
	mockViaCepService := new(MockViaCepService)
	mockTempoService := new(MockTempoService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetCepUseCase(appConfig, mockViaCepService, mockTempoService, tracingProvider)

	cep := "00000-000"
	mockViaCepService.On("GetCep", mock.Anything, cep).Return(entities.ViaCepDto{}, errors.New("cep not found"))

	// Act
	_, err := useCase.GetTempoPorCep(context.Background(), cep)

	// Assert
	require.Error(t, err)
	require.IsType(t, &entities.CustomErrors{}, err)
	customErr := err.(*entities.CustomErrors)
	require.Equal(t, 404, customErr.Code)
	require.Equal(t, "can not find zipcode", customErr.Message)
	mockViaCepService.AssertCalled(t, "GetCep", mock.Anything, cep)
}
