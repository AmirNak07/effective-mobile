package http

import (
	"effective-mobile/internal/http/dto"
	"effective-mobile/internal/service"
	"effective-mobile/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler struct {
	validator *validator.Validate
	logger    logger.Logger

	service service.Subscription
}

func NewHandler(l logger.Logger, s service.Subscription) *Handler {
	return &Handler{
		validator: validator.New(),
		logger:    l,
		service:   s,
	}
}

func (h *Handler) respondWithError(w http.ResponseWriter, r *http.Request, statusCode int, clientMsg string, internalErr error) {
	if internalErr != nil {
		h.logger.Error(r.Context(), "request failed",
			zap.Int("status_code", statusCode),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Error(internalErr),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := dto.ErrorResponse{
		Error:      http.StatusText(statusCode),
		Message:    clientMsg,
		StatusCode: statusCode,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func (h *Handler) respondJSON(w http.ResponseWriter, r *http.Request, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error(r.Context(), "failed to encode response",
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
	}
}
