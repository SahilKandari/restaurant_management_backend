package models

import (
	"time"
)

type Table struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name" validate:"required,max=100,min=1"`
	Capacity     int       `json:"capacity" validate:"required,min=1"`
	RestaurantID int       `json:"restaurant_id" validate:"required"`
	Location     string    `json:"location"`
	Status       string    `json:"status" validate:"required,oneof=available occupied reserved"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
