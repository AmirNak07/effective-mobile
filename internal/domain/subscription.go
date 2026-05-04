package domain

import (
	"time"
)

type Subscription struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"` // Format: MM-YYYY
	EndDate     *string   `json:"end_date"`   // Format: MM-YYYY, optional
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
