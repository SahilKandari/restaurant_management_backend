package models

import (
	"time"
)

type Menu struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name" validate:"required,oneof=appetizer main_course dessert beverage"`
	RestaurantID uint      `json:"restaurant_id" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
