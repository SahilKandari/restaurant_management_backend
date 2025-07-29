package models

import (
	"time"
)

type Invoice struct {
	ID            uint      `json:"id"`
	OrderID       uint      `json:"order_id" validate:"required"`
	Amount        float64   `json:"amount" validate:"required"`
	Tax           float64   `json:"tax"`
	Total         float64   `json:"total"`
	Status        string    `json:"status" validate:"required,oneof=pending paid"`
	PaymentMethod string    `json:"payment_method" validate:"required,oneof=cash credit_card debit_card online"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	RestaurantID  uint      `json:"restaurant_id" validate:"required"`
}
