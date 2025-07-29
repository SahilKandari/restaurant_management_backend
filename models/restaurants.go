package models

import "time"

type Restaurant struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" validate="required,min=3"`
	OwnerID     uint      `json:"owner_id" validate="required"`
	Logo        string    `json:"logo"`
	Address     string    `json:"address" validate="required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
