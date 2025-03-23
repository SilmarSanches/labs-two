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

type MockConsultaService struct {
	mock.Mock
}

func (m *MockConsultaService) GetTempo(ctx context.Context, cep entities.CepRequestDto) (entities.TempoResponseDto, error) {
	args := m.Called(ctx, cep)
	return args.Get(0).(entities.TempoResponseDto), args.Error(1)
}

func TestGet_Success(t *testing.T) {
	// Arrange
	mockConsultaService := new(MockConsultaService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetConsultaUseCase(appConfig, mockConsultaService, tracingProvider)

	cep := entities.CepRequestDto{Cep: "01001000"}
	cidade := "São Paulo"
	mockConsultaService.On("GetTempo", mock.Anything, cep).Return(entities.TempoResponseDto{
		City:       cidade,
		Kelvin:     273.15,
		Celsius:    0,
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
	mockConsultaService.AssertCalled(t, "GetTempo", mock.Anything, cep)
}

func TestGetTempo_CepNotFount(t *testing.T) {
	// Arrange
	mockConsultaService := new(MockConsultaService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetConsultaUseCase(appConfig, mockConsultaService, tracingProvider)

	cep := entities.CepRequestDto{Cep: "00000000"}
	mockConsultaService.On("GetTempo", mock.Anything, cep).Return(entities.TempoResponseDto{}, errors.New("cep not found"))

	// Act
	_, err := useCase.GetTempoPorCep(context.Background(), cep)

	// Assert
	require.Error(t, err)
	require.IsType(t, &entities.CustomErrors{}, err)
	customErr := err.(*entities.CustomErrors)
	require.Equal(t, 404, customErr.Code)
	require.Equal(t, "can not find zipcode", customErr.Message)
	mockConsultaService.AssertCalled(t, "GetTempo", mock.Anything, cep)
}
