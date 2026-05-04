package http

import "time"

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" validate:"required,min=1,max=100"`
	Price       int     `json:"price" validate:"required,gte=0"`
	UserID      string  `json:"user_id" validate:"required,uuid"`
	StartDate   string  `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" validate:"required,min=1,max=100"`
	Price       int     `json:"price" validate:"required,gte=0"`
	UserID      string  `json:"user_id" validate:"required,uuid"`
	StartDate   string  `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
}

type SubscriptionResponse struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListSubscriptionsResponse struct {
	Items  []SubscriptionResponse `json:"items"`
	Total  int                    `json:"total"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}

type GetTotalCostRequest struct {
	UserID      string `json:"user_id" validate:"omitempty,uuid"`
	ServiceName string `json:"service_name" validate:"omitempty,min=1,max=100"`
	From        string `json:"from" validate:"required,datetime=01-2006"`
	To          string `json:"to" validate:"required,datetime=01-2006"`
}

type ErrorResponse struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
