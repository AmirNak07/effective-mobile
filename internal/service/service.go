package service

import (
	"context"
	"effective-mobile/internal/domain"
	"effective-mobile/internal/http/dto"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid" // nolint: golint
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type Repository interface {
	Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error)
	List(ctx context.Context, input dto.ListSubscriptionsRequest) ([]domain.Subscription, int, error)
	Update(ctx context.Context, id uuid.UUID, sub domain.Subscription) (domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetForCost(ctx context.Context, userID string, serviceName string, from string, to string) ([]domain.Subscription, error)
}

type Subscription interface {
	Create(ctx context.Context, input dto.CreateSubscriptionRequest) (domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error)
	List(ctx context.Context, input dto.ListSubscriptionsRequest) ([]domain.Subscription, int, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateSubscriptionRequest) (domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, input dto.GetTotalCostRequest) (int, error)
}

type subscription struct {
	repo Repository
}

func NewService(repo Repository) Subscription {
	return &subscription{
		repo: repo,
	}
}

func (s *subscription) Create(ctx context.Context, input dto.CreateSubscriptionRequest) (domain.Subscription, error) {
	sub := input.ToDomain()
	sub.ID = uuid.New().String()
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()

	return s.repo.Create(ctx, sub)
}

func (s *subscription) Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	return s.repo.Get(ctx, id)
}

func (s *subscription) List(ctx context.Context, input dto.ListSubscriptionsRequest) ([]domain.Subscription, int, error) {
	return s.repo.List(ctx, input)
}

func (s *subscription) Update(ctx context.Context, id uuid.UUID, input dto.UpdateSubscriptionRequest) (domain.Subscription, error) {
	sub := domain.Subscription{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		UpdatedAt:   time.Now(),
	}

	return s.repo.Update(ctx, id, sub)
}

func (s *subscription) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscription) GetTotalCost(ctx context.Context, input dto.GetTotalCostRequest) (int, error) {
	from, err := time.Parse("01-2006", input.From)
	if err != nil {
		return 0, fmt.Errorf("invalid from date: %w", err)
	}

	to, err := time.Parse("01-2006", input.To)
	if err != nil {
		return 0, fmt.Errorf("invalid to date: %w", err)
	}

	if from.After(to) {
		return 0, errors.New("from date must be before or equal to to date")
	}

	subs, err := s.repo.GetForCost(ctx, input.UserID, input.ServiceName, input.From, input.To)
	if err != nil {
		return 0, err
	}

	totalCost := 0
	for _, sub := range subs {
		subStart, err := time.Parse("01-2006", sub.StartDate)
		if err != nil {
			continue // skip invalid data
		}

		var subEnd time.Time
		if sub.EndDate != nil && *sub.EndDate != "" {
			subEnd, err = time.Parse("01-2006", *sub.EndDate)
			if err != nil {
				continue // skip invalid data
			}
		}

		months := calculateActiveMonths(subStart, subEnd, from, to)
		totalCost += months * sub.Price
	}

	return totalCost, nil
}

func calculateActiveMonths(start, end, from, to time.Time) int {
	eStart := start
	if from.After(eStart) {
		eStart = from
	}

	eEnd := to
	if !end.IsZero() && end.Before(eEnd) {
		eEnd = end
	}

	if eStart.After(eEnd) {
		return 0
	}

	years := eEnd.Year() - eStart.Year()
	months := int(eEnd.Month()) - int(eStart.Month())
	return years*12 + months + 1
}
