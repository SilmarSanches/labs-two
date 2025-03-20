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

type ServiceTempoInterface interface {
	GetTempo(ctx context.Context, location string) (entities.TempoResponseDto, error)
}

type ServiceTempo struct {
	HttpClient HttpClient
	appConfig  *config.AppSettings
}

func NewServiceTempo(httpClient HttpClient, appConfig *config.AppSettings) *ServiceTempo {
	return &ServiceTempo{
		HttpClient: httpClient,
		appConfig:  appConfig,
	}
}

func (s *ServiceTempo) GetTempo(ctx context.Context, location string) (entities.TempoResponseDto, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/consulta-tempo", s.appConfig.UrlTempo)
	log.Println("url:", url)

	tempo := entities.TempoRequestDto{
		Location: location,
	}

	tempoJSON, err := json.Marshal(tempo)
	if err != nil {
		log.Println("erro ao serializar objeto tempo:", err)
		return entities.TempoResponseDto{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(tempoJSON))
	if err != nil {
		log.Println("erro ao criar requisição HTTP:", err)
		return entities.TempoResponseDto{}, err
	}

	res, err := s.HttpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("timeout de 5s excedido ao consultar o serviço ConsultaTempo:", err)
			return entities.TempoResponseDto{}, fmt.Errorf("timeout de 5s excedido ao consultar o serviço ConsultaTempo: %v", err)
		}
	}

	if res == nil {
		log.Println("resposta nula ao consultar o tempo")
		if err != nil {
			log.Println("erro ao consultar o serviço tempo:", err)
		}
		return entities.TempoResponseDto{}, errors.New("resposta nula ao consultar o tempo")
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Printf("erro ao fechar o corpo da resposta consultaTempo: %v", err)
			}
		}(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		log.Println("erro ao consultar o serviço tempo:", res.StatusCode)
		return entities.TempoResponseDto{}, fmt.Errorf("erro ao consultar o serviço tempo: %d", res.StatusCode)
	}

	var data entities.TempoResponseDto
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println("erro ao decodificar resposta JSON do servicoTempo:", err)
		return entities.TempoResponseDto{}, fmt.Errorf("erro ao decodificar resposta JSON do servicoTempo: %w", err)
	}

	if data.City == "" {
		log.Println("Location não encontrado:", location)
		return entities.TempoResponseDto{}, fmt.Errorf("location não encontrado: %s", location)
	}

	return data, nil
}
