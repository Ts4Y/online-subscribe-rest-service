package router

import (
	"net/http"
	_ "online-subscribe-rest-service/docs"
	"online-subscribe-rest-service/internal/api/handler"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/subscriptions/{user_id}/list", h.SubscriptionsList)
	r.Get("/subscriptions/{id}", h.SubscriptionByID)
	r.Post("/subscriptions", h.CreateSubscription)
	r.Put("/subscriptions", h.UpdateSubscription)
	r.Delete("/subscriptions/{id}", h.DeleteSubscription)
	r.Get("/subscriptions/sum", h.SubscriptionsSum)
	r.Get("/swagger/*", httpSwagger.Handler())
	return r
}
