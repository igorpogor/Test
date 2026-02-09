package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a subscription record in the database
type Subscription struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServiceName string    `json:"service_name" gorm:"not null"`
	Price       int       `json:"price" gorm:"not null"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	StartDate   string    `json:"start_date" gorm:"not null"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SubscriptionRequest represents the request body for creating/updating a subscription
// @Description Request body for creating or updating a subscription
type SubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int     `json:"price" binding:"required,min=1" example:"400"`
	UserID      string  `json:"user_id" binding:"required,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string  `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2025"`
}

// SubscriptionResponse represents the response body for a subscription
// @Description Response body containing subscription data
type SubscriptionResponse struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price" example:"400"`
	UserID      uuid.UUID `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" example:"07-2025"`
	EndDate     *string   `json:"end_date,omitempty" example:"12-2025"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AggregateRequest represents query parameters for the aggregation endpoint
// @Description Query parameters for aggregating subscription costs
type AggregateRequest struct {
	UserID      *string `json:"user_id,omitempty" form:"user_id"`
	ServiceName *string `json:"service_name,omitempty" form:"service_name"`
	StartDate   *string `json:"start_date,omitempty" form:"start_date"`
	EndDate     *string `json:"end_date,omitempty" form:"end_date"`
}

// AggregateResponse represents the response body for the aggregation endpoint
// @Description Response body containing the total subscription cost
type AggregateResponse struct {
	TotalPrice int `json:"total_price" example:"1200"`
}

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
