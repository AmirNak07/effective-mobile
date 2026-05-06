package http

import (
	"effective-mobile/internal/domain"
	"effective-mobile/internal/http/dto"
	"effective-mobile/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Validation failed", err)
		return
	}

	sub, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.respondWithError(w, r, http.StatusInternalServerError, "Failed to create subscription", err)
		return
	}

	h.respondJSON(w, r, http.StatusCreated, toSubscriptionResponse(sub))
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Invalid subscription ID format", nil)
		return
	}

	sub, err := h.service.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			h.respondWithError(w, r, http.StatusNotFound, "Subscription not found", err)
			return
		}
		h.respondWithError(w, r, http.StatusInternalServerError, "Failed to get subscription", err)
		return
	}

	h.respondJSON(w, r, http.StatusOK, toSubscriptionResponse(sub))
}

func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limitStr := query.Get("limit")
	limit := 20
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.respondWithError(w, r, http.StatusBadRequest, "Invalid limit format", err)
			return
		}
	}

	offsetStr := query.Get("offset")
	offset := 0
	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.respondWithError(w, r, http.StatusBadRequest, "Invalid offset format", err)
			return
		}
	}

	req := dto.ListSubscriptionsRequest{
		UserID: query.Get("user_id"),
		Limit:  limit,
		Offset: offset,
	}

	if err := h.validator.Struct(req); err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Validation failed", err)
		return
	}

	subs, total, err := h.service.List(r.Context(), req)
	if err != nil {
		h.respondWithError(w, r, http.StatusInternalServerError, "Failed to list subscriptions", err)
		return
	}

	items := make([]dto.SubscriptionResponse, len(subs))
	for i, sub := range subs {
		items[i] = toSubscriptionResponse(sub)
	}

	h.respondJSON(w, r, http.StatusOK, dto.ListSubscriptionsResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Invalid subscription ID format", nil)
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if err = h.validator.Struct(req); err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Validation failed", err)
		return
	}

	sub, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			h.respondWithError(w, r, http.StatusNotFound, "Subscription not found", err)
			return
		}
		h.respondWithError(w, r, http.StatusInternalServerError, "Failed to update subscription", err)
		return
	}

	h.respondJSON(w, r, http.StatusOK, toSubscriptionResponse(sub))
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Invalid subscription ID format", nil)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			h.respondWithError(w, r, http.StatusNotFound, "Subscription not found", err)
			return
		}
		h.respondWithError(w, r, http.StatusInternalServerError, "Failed to delete subscription", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	req := dto.GetTotalCostRequest{
		UserID:      query.Get("user_id"),
		ServiceName: query.Get("service_name"),
		From:        query.Get("from"),
		To:          query.Get("to"),
	}

	if err := h.validator.Struct(req); err != nil {
		h.respondWithError(w, r, http.StatusBadRequest, "Validation failed", err)
		return
	}

	totalCost, err := h.service.GetTotalCost(r.Context(), req)
	if err != nil {
		// Business logic errors like "from after to" should be 400 Bad Request
		h.respondWithError(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	h.respondJSON(w, r, http.StatusOK, dto.TotalCostResponse{
		TotalCost: totalCost,
	})
}

func toSubscriptionResponse(sub domain.Subscription) dto.SubscriptionResponse {
	return dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}
}
