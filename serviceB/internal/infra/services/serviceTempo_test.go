package services

import (
	"context"
	"io"
	"labs-two-service-b/config"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetTempo_Success(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlTempo: "https://api.weather.com", TempoApiKey: "fake-key"}
	service := NewServiceTempo(mockHttpClient, appConfig)

	mockResponse := &http.Response{
        StatusCode: http.StatusOK,
        Body: io.NopCloser(strings.NewReader(`{
            "current": {
                "temp_c": 25.5
            }
        }`)),
    }

	mockHttpClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Act
	ctx := context.Background()
	cidade := "São Paulo"
	result, _ := service.GetTempo(ctx, cidade)

	// Assert
	require.Equal(t, 25.5, result.Current.TempC, "Temperatura incorreta")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetTempo_Timeout(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlTempo: "https://api.weather.com", TempoApiKey: "fake-key"}
	service := NewServiceTempo(mockHttpClient, appConfig)

	mockHttpClient.On("Do", mock.Anything).Return((*http.Response)(nil), context.DeadlineExceeded)

	// Act
	ctx := context.Background()
	cidade := "São Paulo"
	_, err := service.GetTempo(ctx, cidade)

	// Assert
	require.ErrorContains(t, err, "timeout de 5s")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetTempo_NilResponse(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlTempo: "https://api.weather.com", TempoApiKey: "fake-key"}
	service := NewServiceTempo(mockHttpClient, appConfig)

	mockHttpClient.On("Do", mock.Anything).Return((*http.Response)(nil), nil)

	// Act
	ctx := context.Background()
	cidade := "São Paulo"
	_, err := service.GetTempo(ctx, cidade)

	// Assert
	require.EqualError(t, err, "resposta nula ao consultar o weatherapi")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetTempo_Non200StatusCode(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlTempo: "https://api.weather.com", TempoApiKey: "fake-key"}
	service := NewServiceTempo(mockHttpClient, appConfig)

	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(``)),
	}

	mockHttpClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Act
	ctx := context.Background()
	cidade := "São Paulo"
	_, err := service.GetTempo(ctx, cidade)

	// Assert
	require.Equal(t, "erro ao consultar o serviço weatherapi: 500", err.Error())
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetTempo_JSONDecodeError(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlTempo: "https://api.weather.com", TempoApiKey: "fake-key"}
	service := NewServiceTempo(mockHttpClient, appConfig)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader("{invalid json}")),
	}

	mockHttpClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Act
	ctx := context.Background()
	cidade := "São Paulo"
	_, err := service.GetTempo(ctx, cidade)

	// Assert
	require.ErrorContains(t, err, "invalid character")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}
