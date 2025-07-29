package models

import (
	"time"
)

type Food struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name" validate:"required,max=100,min=2"`
	Description  string    `json:"description" validate:"max=500"`
	Price        float64   `json:"price" validate:"required"`
	ImageURL     string    `json:"image_url" validate:"max=255"`
	MenuID       uint      `json:"menu_id" validate:"required"`
	RestaurantID uint      `json:"restaurant_id" validate:"required"`
	Ingredients  string    `json:"ingredients"`
	PrepTime     int       `json:"prep_time"` // Preparation time in minutes
	Calories     int       `json:"calories" validate:"min=0"`
	SpicyLevel   int       `json:"spicy_level" validate:"min=0,max=5"`
	Vegetarian   bool      `json:"vegetarian"`
	Available    bool      `json:"available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
