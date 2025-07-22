package service

import (
	"context"
	"fmt"
	"online-subscribe-rest-service/internal/entity"

	"github.com/gofrs/uuid/v5"
)

type Repo interface {
	SubscriptionByID(context.Context, uuid.UUID) (entity.Subscription, error)
	UpdateSubscription(context.Context, entity.Subscription) error
	CreateSubscription(context.Context, entity.Subscription) (uuid.UUID, error)
	DeleteSubscription(context.Context, uuid.UUID) error
	SubscriptionsList(context.Context, uuid.UUID) ([]entity.Subscription, error)
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) UpdateSubscription(ctx context.Context, sub entity.Subscription) error {
	_, err := s.repo.SubscriptionByID(ctx, sub.ID)
	if err != nil {
		return fmt.Errorf("service: failed to find subscription with id %s: %w", sub.ID, err)
	}

	if err := s.repo.UpdateSubscription(ctx, sub); err != nil {
		return fmt.Errorf("service: failed to update subscription: %w", err)
	}

	return nil
}

func (s *Service) CreateSubscription(ctx context.Context, sub entity.Subscription) (uuid.UUID, error) {
	id, err := s.repo.CreateSubscription(ctx, sub)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *Service) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteSubscription(ctx, id); err != nil {
		return fmt.Errorf("failed to delete subscription with id %s: %w", id, err)
	}

	return nil
}

func (s *Service) SubscriptionByID(ctx context.Context, id uuid.UUID) (entity.Subscription, error) {
	sub, err := s.repo.SubscriptionByID(ctx, id)
	if err != nil {
		return entity.Subscription{}, fmt.Errorf("failed to get subscription by id %s: %w", id, err)
	}

	return sub, nil
}

func (s *Service) SubscriptionsList(ctx context.Context, userID uuid.UUID) ([]entity.Subscription, error) {

	subs, err := s.repo.SubscriptionsList(ctx, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions list by userID %s: %w", userID, err)
	}

	return subs, nil

}
