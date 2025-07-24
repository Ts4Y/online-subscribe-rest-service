package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"online-subscribe-rest-service/internal/entity"
	"online-subscribe-rest-service/pkg/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

// @title Subscriptions api docs
// @description REST API for managing subscriptions

type SubscriptionsService interface {
	SubscriptionByID(context.Context, uuid.UUID) (entity.Subscription, error)
	SubscriptionsList(context.Context, uuid.UUID) ([]entity.Subscription, error)
	CreateSubscription(context.Context, entity.Subscription) (uuid.UUID, error)
	UpdateSubscription(context.Context, entity.Subscription) error
	DeleteSubscription(context.Context, uuid.UUID) error
	SubscriptionsSum(ctx context.Context, params entity.SubscriptionsSumParams) (entity.UserSubscriptionsSum, error)
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


// @Summary Get all subscriptions by user_id
// @Description Возвращает список подписок по user_id
// @Tags Subscriptions
// @Param user_id path string true "User ID (UUID)"
// @Success 200 {array} entity.Subscription
// @Failure 400 {string} string "Invalid user_id"
// @Failure 404 {string} string "Subscriptions not found"
// @Failure 500 {string} string "Internal server error"
// @Router       /users/{user_id}/subscriptions [get]
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

// @Summary Get subscription by ID
// @Description Возвращает одну подписку по её ID
// @Tags Subscriptions
// @Param id path string true "Subscription ID (UUID)"
// @Success 200 {object} entity.Subscription
// @Failure 400 {string} string "Invalid subscription ID"
// @Failure 404 {string} string "Subscription not found"
// @Failure 500 {string} string "Internal server error"
// @Router       /subscriptions/{id} [get]
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

// @Summary Create subscription
// @Description Создаёт новую подписку
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body entity.Subscription true "Subscription payload"
// @Success 200 {string} string "Subscription created (ID)"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router       /subscriptions [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var subscription entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		h.log.ErrorF("failed to decode request body to struct: %w", err)
		http.Error(w, "failed to decode request body to struct", http.StatusBadRequest)
	}

	if err := subscription.Validate(); err != nil {
		h.log.ErrorF("handler: incorrect params: %w", err)
		http.Error(w, "handler: incorrect params:", http.StatusBadRequest)
		return
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

// @Summary Update subscription
// @Description Обновляет существующую подписку
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body entity.Subscription true "Subscription payload"
// @Success 200 {object} entity.Subscription
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Subscription not found"
// @Failure 500 {string} string "Internal server error"
// @Router       /subscriptions [put]
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

// @Summary Delete subscription
// @Description Удаляет подписку по её ID
// @Tags Subscriptions
// @Param id path string true "Subscription ID (UUID)"
// @Success 200 {string} string "Subscription successfully deleted"
// @Failure 400 {string} string "Invalid subscription ID"
// @Failure 500 {string} string "Internal server error"
// @Router       /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	qID := chi.URLParam(r, "id")
	id, err := uuid.FromString(qID)
	if err != nil {
		h.log.ErrorF("handler: failed to convert string to int %w", err)
		http.Error(w, "handler: failed to convert string to int", http.StatusBadRequest)
		return
	}

	if err := h.subscriptionsService.DeleteSubscription(ctx, id); err != nil {
		h.log.ErrorF("handler: failed to delete subscription", err)
		http.Error(w, "handler: failed to delete subscription", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("subscription sucessfully deleted"))

	if err != nil {
		h.log.ErrorF("handler: failed to write response %w", err)
		http.Error(w, "handler: failed to write response", http.StatusInternalServerError)
	}
}

// @Summary Get total subscription cost
// @Description Подсчитывает суммарную стоимость подписок за период с фильтрами
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "User ID (UUID)"
// @Param service_name query string true "Service name"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} entity.UserSubscriptionsSum
// @Failure 400 {string} string "Invalid or missing parameters"
// @Failure 500 {string} string "Internal server error"
// @Router       /subscriptions/sum [get]
func (h *Handler) SubscriptionsSum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param, err := parseSubscriptionsSumParams(r.URL.Query())
	if err != nil {
		h.log.ErrorF("handler: failed to get parsed params %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := param.Validate(); err != nil {
		h.log.ErrorF("handler: incorrect params: %w", err)
		http.Error(w, fmt.Sprintf("handler: incorrect params: %s", err.Error()), http.StatusBadRequest)
		return
	}

	subSum, err := h.subscriptionsService.SubscriptionsSum(ctx, param)

	if err != nil {
		h.log.ErrorF("handler: failed to get subscriptions sum %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(subSum); err != nil {
		h.log.ErrorF("handler: failed to encode subscriptions sum %w", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func parseSubscriptionsSumParams(url url.Values) (entity.SubscriptionsSumParams, error) {
	qUserID := url.Get("user_id")
	serviceName := url.Get("service_name")
	qStartDate := url.Get("start_date")
	qEndDate := url.Get("end_date")

	userID, err := uuid.FromString(qUserID)
	if err != nil {
		return entity.SubscriptionsSumParams{}, fmt.Errorf("invalid user_id: %w", err)
	}

	startDate, err := time.Parse(time.DateOnly, qStartDate)
	if err != nil {
		return entity.SubscriptionsSumParams{}, fmt.Errorf("invalid start_date: %w", err)
	}

	param := entity.SubscriptionsSumParams{
		UserID:      userID,
		ServiceName: serviceName,
		StartDate:   startDate,
	}

	if qEndDate != "" {
		endDate, err := time.Parse(time.DateOnly, qEndDate)
		if err != nil {
			return entity.SubscriptionsSumParams{}, fmt.Errorf("invalid end_date: %w", err)
		}

		param.EndDate = &endDate
	}

	return param, nil
}
