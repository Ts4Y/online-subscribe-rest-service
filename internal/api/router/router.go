package router

import (
	"net/http"
	"online-subscribe-rest-service/internal/api/handler"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/subscriptions/{user_id}/list", h.SubscriptionsList)
	r.Get("/subscriptions/{id}", h.SubscriptionByID)
	r.Post("/subscriptions", h.CreateSubscription)
	r.Put("/subscriptions", h.UpdateSubscription)
	r.Delete("/subscriptions/{id}", h.DeleteSubscription)

	return r
}
