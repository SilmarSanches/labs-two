package usecases

import (
	"context"
	"errors"
	"labs-two-serviceb/config"
	"labs-two-serviceb/internal/entities"
	"labs-two-serviceb/internal/infra/tracing"
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

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetTempo(ctx context.Context, city string) (entities.TempoDto, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(entities.TempoDto), args.Error(1)
}

func TestGetTempo_Success(t *testing.T) {
	// Arrange
	mockViaCepService := new(MockViaCepService)
	mockWeatherService := new(MockWeatherService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetTempoUseCase(appConfig, mockViaCepService, mockWeatherService, tracingProvider)

	cep := "01001000"
	cidade := "São Paulo"
	mockViaCepService.On("GetCep", mock.Anything, cep).Return(entities.ViaCepDto{Localidade: cidade}, nil)
	mockWeatherService.On("GetTempo", mock.Anything, cidade).Return(entities.TempoDto{Current: entities.Current{TempC: 25.5}}, nil)

	// Act
	result, err := useCase.GetTempo(context.Background(),cep)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 25.5, result.Celsius)
	require.Equal(t, 298.5, result.Kelvin)
	require.Equal(t, 77.9, result.Fahrenheit)
	mockViaCepService.AssertCalled(t, "GetCep", mock.Anything, cep)
	mockWeatherService.AssertCalled(t, "GetTempo", mock.Anything, cidade)
}

func TestGetTempo_InvalidCep(t *testing.T) {
	// Arrange
	mockViaCepService := new(MockViaCepService)
	mockWeatherService := new(MockWeatherService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetTempoUseCase(appConfig, mockViaCepService, mockWeatherService, tracingProvider)

	cep := "invalid-cep"
	mockViaCepService.On("GetCep", mock.Anything, cep).Return(nil, &entities.CustomError{Code: 422, Message: "invalid zipcode"})

	// Act
	_, err := useCase.GetTempo(context.Background(), cep)

	// Assert
	require.Error(t, err)
	require.IsType(t, &entities.CustomError{}, err)
	customErr := err.(*entities.CustomError)
	require.Equal(t, 422, customErr.Code)
	require.Equal(t, "invalid zipcode", customErr.Message)
}

func TestGetTempo_CepNotFound(t *testing.T) {
	// Arrange
	mockViaCepService := new(MockViaCepService)
	mockWeatherService := new(MockWeatherService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetTempoUseCase(appConfig, mockViaCepService, mockWeatherService, tracingProvider)

	cep := "01001000"
	mockViaCepService.On("GetCep", mock.Anything, cep).Return(entities.ViaCepDto{}, errors.New("CEP not found"))

	// Act
	_, err := useCase.GetTempo(context.Background(), cep)

	// Assert
	require.Error(t, err)
	require.IsType(t, &entities.CustomError{}, err)
	customErr := err.(*entities.CustomError)
	require.Equal(t, 404, customErr.Code)
	require.Equal(t, "can not find zipcode", customErr.Message)
	mockViaCepService.AssertCalled(t, "GetCep", mock.Anything, cep)
}

func TestGetTempo_TemperatureNotFound(t *testing.T) {
	// Arrange
	mockViaCepService := new(MockViaCepService)
	mockWeatherService := new(MockWeatherService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetTempoUseCase(appConfig, mockViaCepService, mockWeatherService, tracingProvider)

	cep := "01001000"
	cidade := "São Paulo"
	mockViaCepService.On("GetCep", mock.Anything, cep).Return(entities.ViaCepDto{Localidade: cidade}, nil)
	mockWeatherService.On("GetTempo", mock.Anything, cidade).Return(entities.TempoDto{}, errors.New("Temperature not found"))

	// Act
	_, err := useCase.GetTempo(context.Background(), cep)

	// Assert
	require.Error(t, err)
	require.IsType(t, &entities.CustomError{}, err)
	customErr := err.(*entities.CustomError)
	require.Equal(t, 404, customErr.Code)
	require.Equal(t, "can not find temperature", customErr.Message)
	mockViaCepService.AssertCalled(t, "GetCep", mock.Anything, cep)
	mockWeatherService.AssertCalled(t, "GetTempo", mock.Anything, cidade)
}
