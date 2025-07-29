package models

import (
	"time"
)

type Note struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" validate:"required,min=3,max=100"`
	Content      string    `json:"content" validate:"required,min=10,max=500"`
	Priority     string    `json:"priority" validate:"required,oneof=low medium high"`
	RestaurantID uint      `json:"restaurant_id" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
