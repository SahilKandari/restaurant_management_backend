package models

import (
	"time"
)

type Order struct {
	ID           uint        `json:"id"`
	TableID      uint        `json:"table_id" validate:"required"`
	RestaurantID uint        `json:"restaurant_id" validate:"required"`
	OrderDate    time.Time   `json:"order_date" validate:"required"`
	TotalPrice   float64     `json:"total_price"`
	Status       string      `json:"status" validate:"required,oneof=pending preparing ready served paid cancelled"`
	Notes        string      `json:"notes"`
	OrderItems   []OrderItem `json:"order_items"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type OrderStatus struct {
	Status string `json:"status" validate:"required,oneof=pending preparing ready served paid cancelled"`
}
