package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username" validate:"required,min=3,max=20"`
	Password  string    `json:"password,omitempty" validate:"required,min=6,max=100"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone" validate:"required,min=10,max=15"`
	Role      string    `json:"role" validate:"required,oneof=Admin User"`
	Token     string    `json:"token,omitempty"`
	AvatarURL string    `json:"avatar_url" validate:"omitempty,url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserWithOldPassword struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username" validate:"required,min=3,max=20"`
	Password    string    `json:"password,omitempty" validate:"required,min=6,max=100"`
	Email       string    `json:"email" validate:"required,email"`
	Phone       string    `json:"phone" validate:"required,min=10,max=15"`
	Role        string    `json:"role" validate:"required,oneof=Admin User"`
	Token       string    `json:"token,omitempty"`
	AvatarURL   string    `json:"avatar_url" validate:"omitempty,url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OldPassword string    `json:"old_password,omitempty" validate:"required,min=6,max=100"`
}

type ConfirmPassword struct {
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type PasswordResetEmail struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Email       string `json:"email" validate:"required,email"`
	OTP         string `json:"otp" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
}
