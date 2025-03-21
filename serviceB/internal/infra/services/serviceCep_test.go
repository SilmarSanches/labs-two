package services

import (
	"context"
	"io"
	"labs-two-serviceb/config"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCep_Success(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlCep: "https://viacep.com.br/ws"}
	service := NewServiceCep(mockHttpClient, appConfig)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(`{
			"cep": "01001-000",
			"logradouro": "Praça da Sé",
			"complemento": "lado ímpar",
			"bairro": "Sé",
			"localidade": "São Paulo",
			"uf": "SP",
			"ibge": "3550308",
			"gia": "1004",
			"ddd": "11",
			"siafi": "7107"
		}`)),
	}

	mockHttpClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Act
	ctx := context.Background()
	cep := "01001000"
	result, _ := service.GetCep(ctx, cep)

	// Assert
	require.Equal(t, "01001-000", result.Cep, "O CEP retornado deve ser 01001-000")
	require.Equal(t, "Praça da Sé", result.Logradouro, "Logradouro incorreto")
	require.Equal(t, "lado ímpar", result.Complemento, "Complemento incorreto")
	require.Equal(t, "Sé", result.Bairro, "Bairro incorreto")
	require.Equal(t, "São Paulo", result.Localidade, "Cidade incorreta")
	require.Equal(t, "SP", result.Uf, "UF incorreta")
	require.Equal(t, "3550308", result.Ibge, "Código IBGE incorreto")
	require.Equal(t, "1004", result.Gia, "GIA incorreta")
	require.Equal(t, "11", result.Ddd, "DDD incorreto")
	require.Equal(t, "7107", result.Siafi, "Código SIAFI incorreto")

	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetCep_Timeout(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlCep: "https://viacep.com.br/ws"}
	service := NewServiceCep(mockHttpClient, appConfig)

	mockHttpClient.On("Do", mock.Anything).Return((*http.Response)(nil), context.DeadlineExceeded)

	// Act
	ctx := context.Background()
	cep := "01001000"
	_, err := service.GetCep(ctx, cep)

	// Assert
	require.ErrorContains(t, err, "timeout de 5s")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetCep_NilResponse(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlCep: "http://viacep.com.br/ws"}
	service := &ServiceCep{HttpClient: mockHttpClient, appConfig: appConfig}

	mockHttpClient.On("Do", mock.Anything).Return((*http.Response)(nil), nil)

	// Act
	ctx := context.Background()
	cep := "01001000"
	_, err := service.GetCep(ctx, cep)

	// Assert
	require.EqualError(t, err, "resposta nula ao consultar o viacep")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetCep_Non200StatusCode(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlCep: "http://viacep.com.br/ws"}
	service := &ServiceCep{HttpClient: mockHttpClient, appConfig: appConfig}

	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(``)),
	}

	mockHttpClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Act
	ctx := context.Background()
	cep := "01001000"
	_, err := service.GetCep(ctx, cep)

	// Assert
	require.Equal(t, "erro ao consultar o serviço ViaCep: 500", err.Error())
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}

func TestGetCep_JSONDecodeError(t *testing.T) {
	// Arrange
	mockHttpClient := new(MockHttpClient)
	appConfig := &config.AppSettings{UrlCep: "http://viacep.com.br/ws"}
	service := &ServiceCep{HttpClient: mockHttpClient, appConfig: appConfig}

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("{invalid json}")),
	}

	mockHttpClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Act
	ctx := context.Background()
	cep := "01001000"
	_, err := service.GetCep(ctx, cep)

	// Assert
	require.Contains(t, err.Error(), "erro ao decodificar resposta JSON do ViaCep")
	mockHttpClient.AssertCalled(t, "Do", mock.Anything)
}
