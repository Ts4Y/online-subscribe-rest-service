package entity

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

func (s Subscription) Validate() error {
	if s.ServiceName == "" {
		return errors.New("service name is empty")
	}

	if s.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if s.StartDate.IsZero() {
		return errors.New("start date is empty")
	}

	if s.EndDate != nil {
		if s.EndDate.IsZero() {
			return errors.New("end date is empty")
		}

		if s.EndDate.Before(s.StartDate) {
			return errors.New("end date must be greater than start date")
		}
	}

	return nil
}

type SubscriptionsSumParams struct {
	UserID      uuid.UUID
	ServiceName string
	StartDate   time.Time
	EndDate     *time.Time
}

func (s SubscriptionsSumParams) Validate() error {
	if s.ServiceName == "" {
		return errors.New("service name is empty")
	}

	if s.StartDate.IsZero() {
		return errors.New("start date must be before end date")
	}

	if s.EndDate != nil {
		if s.EndDate.IsZero() {
			return errors.New("end date is empty")
		}

		if s.EndDate.Before(s.StartDate) {
			return errors.New("end date must be greater than start date")
		}
	}

	return nil
}

type UserSubscriptionsSum struct {
	UserID     uuid.UUID `json:"user_id"`
	TotalPrice int       `json:"total_price"`
}
