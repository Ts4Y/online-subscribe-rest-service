package repository

import (
	"context"
	"errors"
	"fmt"
	"online-subscribe-rest-service/internal/entity"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

type SubscriptionRepo struct {
	db *pgx.Conn
}

func NewSubscriptionRepo(db *pgx.Conn) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) CreateSubscription(ctx context.Context, s entity.Subscription) (uuid.UUID, error) {

	id := uuid.Must(uuid.NewV4())

	query := `
	INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, id, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate)

	if err != nil {
		return uuid.Nil, fmt.Errorf("repository: create subscription: %w", err)
	}

	return id, nil
}

func (r *SubscriptionRepo) UpdateSubscription(ctx context.Context, s entity.Subscription) error {

	query := `
	UPDATE subscriptions
	SET 
	service_name = $1,
	price = $2,
	user_id = $3,
	start_date = $4,
	end_date = $5
	WHERE id = $6
	`

	_, err := r.db.Exec(ctx, query, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate, s.ID)

	if err != nil {
		return fmt.Errorf("repository: update subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepo) DeleteSubscription(ctx context.Context, id uuid.UUID) error {

	query := `
DELETE FROM subscriptions 
WHERE id = $1
`

	_, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("repository: delete subscription: %w", err)
	}

	return nil
}

func (r *SubscriptionRepo) SubscriptionByID(ctx context.Context, id uuid.UUID) (entity.Subscription, error) {

	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions 
	WHERE id = $1
	`

	var subscription entity.Subscription
	err := r.db.QueryRow(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Subscription{}, entity.ErrNotFound
		}

		return entity.Subscription{}, fmt.Errorf("repository: SubscriptionByID: %w", err)
	}

	return subscription, nil
}

func (r *SubscriptionRepo) SubscriptionsList(ctx context.Context, userID uuid.UUID) ([]entity.Subscription, error) {

	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions 
	WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("repository: SubscriptionList: %w", err)
	}

	defer rows.Close()

	var subscriptions []entity.Subscription
	for rows.Next() {
		var s entity.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate); err != nil {

			return nil, fmt.Errorf("repository: SubscriptionList: rows.Scan() %w", err)
		}

		subscriptions = append(subscriptions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: SubscriptionList: rows.Err() %w", err)
	}

	return subscriptions, nil

}
