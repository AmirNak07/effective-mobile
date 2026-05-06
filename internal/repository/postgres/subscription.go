package postgres

import (
	"context"
	"effective-mobile/internal/domain"
	"effective-mobile/internal/http/dto"
	"effective-mobile/internal/service"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) service.Repository {
	return &subscriptionRepository{
		pool: pool,
	}
}

func (r *subscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var created domain.Subscription
	err := r.pool.QueryRow(ctx, query,
		sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.CreatedAt, sub.UpdatedAt,
	).Scan(
		&created.ID, &created.ServiceName, &created.Price, &created.UserID, &created.StartDate, &created.EndDate, &created.CreatedAt, &created.UpdatedAt,
	)

	if err != nil {
		return domain.Subscription{}, fmt.Errorf("failed to create subscription: %w", err)
	}

	return created, nil
}

func (r *subscriptionRepository) Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var sub domain.Subscription
	err := r.pool.QueryRow(ctx, query, id.String()).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Subscription{}, service.ErrSubscriptionNotFound
		}
		return domain.Subscription{}, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub, nil
}

func (r *subscriptionRepository) List(ctx context.Context, input dto.ListSubscriptionsRequest) ([]domain.Subscription, int, error) {
	var whereClauses []string
	var args []any
	argID := 1

	if input.UserID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("user_id = $%d", argID))
		args = append(args, input.UserID)
		argID++
	}

	where := ""
	if len(whereClauses) > 0 {
		where = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subscriptions %s", where)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := input.Offset

	query := fmt.Sprintf(`
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argID, argID+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		err := rows.Scan(
			&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}

	return subs, total, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, id uuid.UUID, sub domain.Subscription) (domain.Subscription, error) {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var updated domain.Subscription
	err := r.pool.QueryRow(ctx, query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.UpdatedAt, id.String(),
	).Scan(
		&updated.ID, &updated.ServiceName, &updated.Price, &updated.UserID, &updated.StartDate, &updated.EndDate, &updated.CreatedAt, &updated.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Subscription{}, service.ErrSubscriptionNotFound
		}
		return domain.Subscription{}, fmt.Errorf("failed to update subscription: %w", err)
	}

	return updated, nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	res, err := r.pool.Exec(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if res.RowsAffected() == 0 {
		return service.ErrSubscriptionNotFound
	}

	return nil
}

func (r *subscriptionRepository) GetForCost(ctx context.Context, userID string, serviceName string, from string, to string) ([]domain.Subscription, error) {
	var whereClauses []string
	var args []any
	argID := 1

	if userID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("user_id = $%d", argID))
		args = append(args, userID)
		argID++
	}

	if serviceName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("service_name = $%d", argID))
		args = append(args, serviceName)
		argID++
	}

	// SQL-level date filtering using to_date for MM-YYYY format
	// Subscription is active during [from, to] if:
	// start_date <= to AND (end_date IS NULL OR end_date >= from)
	whereClauses = append(whereClauses, fmt.Sprintf("to_date(start_date, 'MM-YYYY') <= to_date($%d, 'MM-YYYY')", argID))
	args = append(args, to)
	argID++

	whereClauses = append(whereClauses, fmt.Sprintf("(end_date IS NULL OR end_date = '' OR to_date(end_date, 'MM-YYYY') >= to_date($%d, 'MM-YYYY'))", argID))
	args = append(args, from)

	where := ""
	if len(whereClauses) > 0 {
		where = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		%s
	`, where)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions for cost: %w", err)
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		err := rows.Scan(
			&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}

	return subs, nil
}
