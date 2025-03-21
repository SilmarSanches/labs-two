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
	"net/url"
	"time"
)

type ServiceTempoInterface interface {
	GetTempo(ctx context.Context, cidade string) (entities.TempoDto, error)
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

func (s *ServiceTempo) GetTempo(ctx context.Context, cidade string) (entities.TempoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cidadeEncoded := url.QueryEscape(cidade)
	url := fmt.Sprintf("%s/current.json?q=%s&key=%s", s.appConfig.UrlTempo, cidadeEncoded, s.appConfig.TempoApiKey)
	log.Println("url:", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println("erro ao criar requisição HTTP:", err)
		return entities.TempoDto{}, err
	}

	res, err := s.HttpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("timeout de 5s excedido ao consultar o serviço ViaCep:", err)
			return entities.TempoDto{}, fmt.Errorf("timeout de 5s excedido ao consultar o serviço ViaCep: %v", err)
		}
	}

	if res == nil {
		log.Println("resposta nula ao consultar o weatherapi")
		return entities.TempoDto{}, errors.New("resposta nula ao consultar o weatherapi")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("erro ao fechar o corpo da resposta weatherapi:", err)
			fmt.Printf("erro ao fechar o corpo da resposta weatherapi: %v", err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		log.Println("erro ao consultar o serviço weatherapi:", res.StatusCode)
		return entities.TempoDto{}, fmt.Errorf("erro ao consultar o serviço weatherapi: %v", res.StatusCode)
	}

	var data entities.TempoDto
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println("erro ao decodificar resposta JSON weatherapi:", err)
		return entities.TempoDto{}, err
	}

	return data, nil
}
