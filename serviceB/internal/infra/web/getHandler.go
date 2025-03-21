package web

import (
	"encoding/json"
	"labs-two-serviceb/config"
	"labs-two-serviceb/internal/entities"
	"labs-two-serviceb/internal/infra/services"
	"labs-two-serviceb/internal/infra/tracing"
	"labs-two-serviceb/internal/usecases"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type GetTempoHandler struct {
	config                *config.AppSettings
	GetTempoUseCase       usecases.GetTempoUseCaseInterface
	ServiceCepInterface   services.ServiceCepInterface
	ServiceTempoInterface services.ServiceTempoInterface
	tracingProvider       *tracing.TracingProvider
}

func NewGetTempoHandler(appConfig *config.AppSettings, getTempoUseCase usecases.GetTempoUseCaseInterface, serviceCep services.ServiceCepInterface, serviceTempo services.ServiceTempoInterface, tracingProvider *tracing.TracingProvider) *GetTempoHandler {
	return &GetTempoHandler{
		config:                appConfig,
		GetTempoUseCase:       getTempoUseCase,
		ServiceCepInterface:   serviceCep,
		ServiceTempoInterface: serviceTempo,
		tracingProvider:       tracingProvider,
	}
}

// HandleLabsOne godoc
// @Summary Consulta temperatura baseado no CEP
// @Description Consulta a temperatura atual baseada no CEP fornecido
// @Tags Labs-Two-ServiceB
// @Accept json
// @Produce json
// @Param request body entities.CepRequestDto true "Consulta temperatura"
// @Success 200 {object} entities.GetTempoResponseDto "OK"
// @Failure 404 {object} entities.CustomError "Not Found"
// @Failure 422 {object} entities.CustomError "Invalid Zipcode"
// @Router /consulta-tempo [post]
func (h *GetTempoHandler) HandleLabsOne(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	ctx, span := h.tracingProvider.Tracer.Start(ctx, "HandleConsultaTempo")
	defer span.End()

	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "cep is required", http.StatusBadRequest)
		return
	}

	response, err := h.GetTempoUseCase.GetTempo(ctx, cep)
	if err != nil {
		customErr, ok := err.(*entities.CustomError)
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
