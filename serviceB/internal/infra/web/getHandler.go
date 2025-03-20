package web

import (
	"encoding/json"
	"labs-two-service-b/config"
	"labs-two-service-b/internal/entities"
	"labs-two-service-b/internal/infra/services"
	"labs-two-service-b/internal/infra/tracing"
	"labs-two-service-b/internal/usecases"
	"net/http"
)

type GetTempoHandler struct {
	config                *config.AppSettings
	GetTempoUseCase       usecases.GetTempoUseCaseInterface
	ServiceTempoInterface services.ServiceTempoInterface
	tracingProvider       *tracing.TracingProvider
}

func NewGetCepHandler(appConfig *config.AppSettings, getTempoUseCase usecases.GetTempoUseCaseInterface, serviceTempo services.ServiceTempoInterface, tracingProvider *tracing.TracingProvider) *GetTempoHandler {
	return &GetTempoHandler{
		config:                appConfig,
		GetTempoUseCase:       getTempoUseCase,
		ServiceTempoInterface: serviceTempo,
		tracingProvider:       tracingProvider,
	}
}

// HandleLabsTwo godoc
// @Summary Consulta Temperatura
// @Description Consulta temperatura por cidade
// @Tags Labs-Two
// @Accept json
// @Produce json
// @Param request body entities.TempoRequestDto true "City Request"
// @Success 200 {object} entities.TempoDto "OK"
// @Failure 404 {object} entities.CustomErrors "Not Found"
// @Failure 422 {object} entities.CustomErrors "Invalid Zipcode"
// @Router /consulta-tempo [post]
func (h *GetTempoHandler) HandleLabsTwo(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracingProvider.Tracer.Start(r.Context(), "Consulta Temperatura")
	defer span.End()

	var req entities.TempoRequestDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.GetTempoUseCase.GetTempo(ctx, req.Location)
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
