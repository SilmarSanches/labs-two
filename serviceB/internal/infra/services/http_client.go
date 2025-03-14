package services

import (
	"net/http"
	"time"
)

// Interface para requisições HTTP
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Adaptador que implementa a interface `HttpClient`
type HttpClientAdapter struct {
	client *http.Client
}

// Implementação do método `Do`
func (h *HttpClientAdapter) Do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}

// Agora `NewHttpClient()` retorna um `HttpClient`, e não `*http.Client`
func NewHttpClient() HttpClient {
	return &HttpClientAdapter{client: &http.Client{Timeout: 5 * time.Second}}
}
