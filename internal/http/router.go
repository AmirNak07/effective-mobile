package http

import "net/http"

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.Health)

	mux.HandleFunc("POST /subscriptions", h.CreateSubscription)
	mux.HandleFunc("GET /subscriptions", h.ListSubscriptions)
	mux.HandleFunc("GET /subscriptions/{id}", h.GetSubscription)
	mux.HandleFunc("GET /subscriptions/total-cost", h.GetTotalCost)
	mux.HandleFunc("PUT /subscriptions/{id}", h.UpdateSubscription)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.DeleteSubscription)

	return CORSMiddleware(mux)
}
