package models

import "time"

type OrderItem struct {
	ID        uint      `json:"id"`
	OrderID   uint      `json:"order_id" validate:"required"`
	FoodID    uint      `json:"food_id" validate:"required"`
	FoodName  string    `json:"food_name"`
	Quantity  uint      `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	SubTotal  float64   `json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateOrderItem struct {
	Quantity uint `json:"quantity" validate:"required"`
}
