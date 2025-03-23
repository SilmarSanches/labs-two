package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"log"
	"net/http"
	"time"
)

type ServiceConsultaInterface interface {
	GetTempo(ctx context.Context, cep entities.CepRequestDto) (entities.TempoResponseDto, error)
}

type ServiceConsulta struct {
	HttpClient HttpClient
	appConfig  *config.AppSettings
}

func NewServiceConsulta(httpClient HttpClient, appConfig *config.AppSettings) *ServiceConsulta {
	return &ServiceConsulta{
		HttpClient: httpClient,
		appConfig:  appConfig,
	}
}

func (s *ServiceConsulta) GetTempo(ctx context.Context, cep entities.CepRequestDto) (entities.TempoResponseDto, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/consulta-tempo", s.appConfig.UrlConsulta)
	log.Println("url:", url)

	body, err := json.Marshal(cep)
	if err != nil {
		return entities.TempoResponseDto{}, fmt.Errorf("erro ao serializar o objeto cep: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return entities.TempoResponseDto{}, err
	}

	res, err := s.HttpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return entities.TempoResponseDto{}, fmt.Errorf("timeout de 5s excedido ao consultar o serviço ViaCep: %v", err)
		}
	}

	if res == nil {
		return entities.TempoResponseDto{}, errors.New("resposta nula ao consultar o viacep")
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Printf("erro ao fechar o corpo da resposta: %v", err)
			}
		}(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return entities.TempoResponseDto{}, fmt.Errorf("erro ao consultar o tempo: %d", res.StatusCode)
	}

	var data entities.TempoResponseDto
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return entities.TempoResponseDto{}, fmt.Errorf("erro ao decodificar resposta JSON: %w", err)
	}

	if data.City == "" {
		return entities.TempoResponseDto{}, fmt.Errorf("City não encontrado: %s", cep)
	}

	return data, nil
}
