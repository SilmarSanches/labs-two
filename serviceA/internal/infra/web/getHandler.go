package web

import (
	"encoding/json"
	"labs-two-service-a/config"
	"labs-two-service-a/internal/entities"
	"labs-two-service-a/internal/infra/services"
	"labs-two-service-a/internal/infra/tracing"
	"labs-two-service-a/internal/usecases"
	"net/http"
)

type GetConsultaHandler struct {
	config                   *config.AppSettings
	GetCepUseCase            usecases.GetConsultaUseCaseInterface
	ServiceConsultaInterface services.ServiceConsultaInterface
	tracingProvider          *tracing.TracingProvider
}

func NewGetConsultaHandler(appConfig *config.AppSettings, getConsultaUseCase usecases.GetConsultaUseCaseInterface, serviceConsulta services.ServiceConsultaInterface, tracingProvider *tracing.TracingProvider) *GetConsultaHandler {
	return &GetConsultaHandler{
		config:                   appConfig,
		GetCepUseCase:            getConsultaUseCase,
		ServiceConsultaInterface: serviceConsulta,
		tracingProvider:          tracingProvider,
	}
}

// HandleLabsTwo godoc
// @Summary Consulta CEP
// @Description Consulta dados do CEP fornecido via JSON no corpo da requisição
// @Tags Labs-Two
// @Accept json
// @Produce json
// @Param request body entities.CepRequestDto true "CEP Request"
// @Success 200 {object} entities.TempoResponseDto "OK"
// @Failure 404 {object} entities.CustomErrors "Not Found"
// @Failure 422 {object} entities.CustomErrors "Invalid Zipcode"
// @Router /consulta-cep [post]
func (h *GetConsultaHandler) HandleLabsTwo(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracingProvider.Tracer.Start(r.Context(), "Consulta CEP")
	defer span.End()

	var req entities.CepRequestDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.GetCepUseCase.GetTempoPorCep(ctx, req)
	if err != nil {
		customErr, ok := err.(*entities.CustomErrors)
		if ok {
			w.WriteHeader(customErr.Code)
			json.NewEncoder(w).Encode(customErr)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
