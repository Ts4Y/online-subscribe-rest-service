package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"online-subscribe-rest-service/internal/entity"
	"online-subscribe-rest-service/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

type SubscriptionsService interface {
	SubscriptionByID(context.Context, uuid.UUID) (entity.Subscription, error)
	SubscriptionsList(context.Context, uuid.UUID) ([]entity.Subscription, error)
	CreateSubscription(context.Context, entity.Subscription) (uuid.UUID, error)
	UpdateSubscription(context.Context, entity.Subscription) error
	DeleteSubscription(context.Context, uuid.UUID) error
}

type Handler struct {
	log                  logger.Logger
	subscriptionsService SubscriptionsService
}

func NewHandler(log logger.Logger, subscriptionsService SubscriptionsService) *Handler {
	return &Handler{
		log:                  log,
		subscriptionsService: subscriptionsService,
	}
}

func (h *Handler) SubscriptionsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qUserID := chi.URLParam(r, "user_id")

	userID, err := uuid.FromString(qUserID)
	if err != nil {
		h.log.ErrorF("invalid user_id %s: %w", qUserID, err)
		http.Error(w, fmt.Sprintf("invalid user_id: %s", qUserID), http.StatusBadRequest)
		return
	}

	subscriptions, err := h.subscriptionsService.SubscriptionsList(ctx, userID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, fmt.Sprintf("subscriptions for user_id %s not found", userID), http.StatusNotFound)
			return
		}

		h.log.ErrorF("handler: failed to get subscriptions for user_id %s: %v", userID, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(subscriptions); err != nil {
		h.log.ErrorF("handler: failed to encode subscriptions for user_id %s: %v", userID, err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) SubscriptionByID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	qID := chi.URLParam(r, "id")
	id, err := uuid.FromString(qID)

	if err != nil {
		h.log.ErrorF("invalid id %s: %w", id, err)
		http.Error(w, fmt.Sprintf("subscription for id %s not found", qID), http.StatusBadRequest)
		return
	}

	subscription, err := h.subscriptionsService.SubscriptionByID(ctx, id)

	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, fmt.Sprintf("subscriptions by id %s not found", id), http.StatusNotFound)
			return
		}

		h.log.ErrorF("handler: failed to get subscriptions by id %s: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(subscription); err != nil {
		h.log.ErrorF("handler: failed to encode subscription %w", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var subscription entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		h.log.ErrorF("failed to decode request body to struct: %w", err)
		http.Error(w, "failed to decode request body to struct", http.StatusBadRequest)
	}

	id, err := h.subscriptionsService.CreateSubscription(ctx, subscription)

	if err != nil {
		h.log.ErrorF("handler: failed to create subscription %w", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(id); err != nil {
		h.log.ErrorF("handler: failed to encode id %w", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var subscription entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		h.log.ErrorF("failed to decode request body to struct: %w", err)
		http.Error(w, "handler: failed to decode request body to struct", http.StatusBadRequest)
		return
	}

	if err := h.subscriptionsService.UpdateSubscription(ctx, subscription); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			h.log.ErrorF("handler: failed to update subscription %w", err)
			http.Error(w, "handler: subscription not found", http.StatusNotFound)
		}

		http.Error(w, "handler: failed to update subscription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(subscription); err != nil {
		h.log.ErrorF("handler: failed to encode id %w", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	qID := chi.URLParam(r, "id")
	id, err := uuid.FromString(qID)
	if err != nil {
		http.Error(w, "handler: failed to convert string to int", http.StatusBadRequest)
		return
	}

	if err := h.subscriptionsService.DeleteSubscription(ctx, id); err != nil {

		http.Error(w, "handler: failed to delete subscription", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("subscription sucessfully deleted"))

	if err != nil {
		http.Error(w, "handler: failed to write response", http.StatusInternalServerError)
	}
}
