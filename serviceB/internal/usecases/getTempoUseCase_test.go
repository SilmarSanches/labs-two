package usecases

import (
	"context"
	"errors"
	"labs-two-service-b/config"
	"labs-two-service-b/internal/entities"
	"labs-two-service-b/internal/infra/tracing"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTempoService struct {
	mock.Mock
}

func (m *MockTempoService) GetTempo(ctx context.Context, location string) (entities.TempoDto, error) {
	args := m.Called(ctx, location)
	return args.Get(0).(entities.TempoDto), args.Error(1)
}

func TestGetCep_Success(t *testing.T) {
	// Arrange
	mocktempoService := new(MockTempoService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetTempoUseCase(appConfig, mocktempoService, tracingProvider)

	cidade := "São Paulo"
	mocktempoService.On("GetTempo", mock.Anything, cidade).Return(entities.TempoDto{Location: entities.Location{Name: cidade}}, nil)

	// Act
	result, err := useCase.GetTempo(context.Background(), cidade)

	// Assert
	require.NoError(t, err)
	require.Equal(t, cidade, result.City)
	mocktempoService.AssertCalled(t, "GetTempo", mock.Anything, cidade)
}

func TestGetTempo_CityNotFount(t *testing.T) {
	// Arrange
	mockTempoService := new(MockTempoService)
	appConfig := &config.AppSettings{}
	tracingProvider, _, _ := tracing.NewTracingProvider(tracing.TracingConfig{
		ZipkinURL:   "http://localhost:9411/api/v2/spans",
		ServiceName: "test-service",})

	useCase := NewGetTempoUseCase(appConfig, mockTempoService, tracingProvider)

	location := ""
	mockTempoService.On("GetTempo", mock.Anything, location).Return(entities.TempoDto{}, errors.New("City not found"))

	// Act
	_, err := useCase.GetTempo(context.Background(), location)

	// Assert
	require.Error(t, err)
	require.IsType(t, &entities.CustomErrors{}, err)
	customErr := err.(*entities.CustomErrors)
	require.Equal(t, 422, customErr.Code)
	require.Equal(t, "invalid location", customErr.Message)
}