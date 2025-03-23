package services

import (
	"context"
	"io"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCep_Success(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlConsulta: "http://localhost:8081/consulta-tempo"}
	service := NewServiceConsulta(mockHttpClient, appConfig)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(`{
			"temp_K": 273.15,
			"temp_C": 0,
			"temp_F": 32,
			"city": "São Paulo"
		}`)),
	}

	mockHttpClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Act
	ctx := context.Background()
	cep := entities.CepRequestDto{Cep: "01001000"}
	result, _ := service.GetTempo(ctx, cep)

	// Assert
	require.Equal(t, 273.15, result.Kelvin, "Temperatura em Kelvin incorreta")
	require.Equal(t, 0.0, result.Celsius, "Temperatura em Celsius incorreta")
	require.Equal(t, 32.0, result.Fahrenheit, "Temperatura em Fahrenheit incorreta")
	require.Equal(t, "São Paulo", result.City, "Cidade incorreta")

	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetCep_Timeout(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlConsulta: "http://localhost:8081/consulta-tempo"}
	service := NewServiceConsulta(mockHttpClient, appConfig)

	mockHttpClient.On("Do", mock.Anything).Return((*http.Response)(nil), context.DeadlineExceeded)

	// Act
	ctx := context.Background()
	cep := entities.CepRequestDto{Cep: "01001000"}
	_, err := service.GetTempo(ctx, cep)

	// Assert
	require.ErrorContains(t, err, "timeout de 5s")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetCep_NilResponse(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlConsulta: "http://localhost:8081/consulta-tempo"}
	service := NewServiceConsulta(mockHttpClient, appConfig)

	mockHttpClient.On("Do", mock.Anything).Return((*http.Response)(nil), nil)

	// Act
	ctx := context.Background()
	cep := entities.CepRequestDto{Cep: "01001000"}
	_, err := service.GetTempo(ctx, cep)

	// Assert
	require.EqualError(t, err, "resposta nula ao consultar o viacep")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetCep_Non200StatusCode(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlConsulta: "http://localhost:8081/consulta-tempo"}
	service := NewServiceConsulta(mockHttpClient, appConfig)

	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(``)),
	}

	mockHttpClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Act
	ctx := context.Background()
	cep := entities.CepRequestDto{Cep: "01001000"}
	_, err := service.GetTempo(ctx, cep)

	// Assert
	require.Equal(t, "erro ao consultar o tempo: 500", err.Error())
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}
