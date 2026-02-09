package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServiceName string    `json:"service_name" gorm:"not null" validate:"required"`
	Price       int       `json:"price" gorm:"not null" validate:"required,min=0"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null" validate:"required"`
	StartDate   string    `json:"start_date" gorm:"not null" validate:"required,datetime=01-2006"`
	EndDate     *string   `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SubscriptionRequest struct {
	ServiceName string  `json:"service_name" validate:"required"`
	Price       int     `json:"price" validate:"required,min=0"`
	UserID      string  `json:"user_id" validate:"required,uuid"`
	StartDate   string  `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
}

type SubscriptionResponse struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AggregateRequest struct {
	UserID      *string `json:"user_id,omitempty" form:"user_id" validate:"omitempty,uuid"`
	ServiceName *string `json:"service_name,omitempty" form:"service_name"`
	StartDate   *string `json:"start_date,omitempty" form:"start_date" validate:"omitempty,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" form:"end_date" validate:"omitempty,datetime=01-2006"`
}

type AggregateResponse struct {
	TotalPrice int `json:"total_price"`
}
