package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"labs-two-serviceb/config"
	"labs-two-serviceb/internal/entities"
	"log"

	"net/http"
	"time"
)

type ServiceCepInterface interface {
	GetCep(ctx context.Context, cep string) (entities.ViaCepDto, error)
}

type ServiceCep struct {
	HttpClient HttpClient
	appConfig  *config.AppSettings
}

func NewServiceCep(httpClient HttpClient, appConfig *config.AppSettings) *ServiceCep {
	return &ServiceCep{
		HttpClient: httpClient,
		appConfig:  appConfig,
	}
}

func (s *ServiceCep) GetCep(ctx context.Context, cep string) (entities.ViaCepDto, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/%s/json", s.appConfig.UrlCep, cep)
	log.Println("url:", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println("erro ao criar requisição HTTP:", err)
		return entities.ViaCepDto{}, err
	}

	res, err := s.HttpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("timeout de 5s excedido ao consultar o serviço ViaCep:", err)
			return entities.ViaCepDto{}, fmt.Errorf("timeout de 5s excedido ao consultar o serviço ViaCep: %v", err)
		}
	}

	if res == nil {
		log.Println("resposta nula ao consultar o viacep")
		if err != nil {
			log.Println("erro ao consultar o serviço ViaCep:", err)
		}
		return entities.ViaCepDto{}, errors.New("resposta nula ao consultar o viacep")
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Printf("erro ao fechar o corpo da resposta ViaCep: %v", err)
			}
		}(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		log.Println("erro ao consultar o serviço ViaCep:", res.StatusCode)
		return entities.ViaCepDto{}, fmt.Errorf("erro ao consultar o serviço ViaCep: %d", res.StatusCode)
	}

	var data entities.ViaCepDto
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println("erro ao decodificar resposta JSON do ViaCep:", err)
		return entities.ViaCepDto{}, fmt.Errorf("erro ao decodificar resposta JSON do ViaCep: %w", err)
	}

	return data, nil
}
