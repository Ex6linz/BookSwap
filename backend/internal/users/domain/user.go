package domain

import (
	"time"

	"github.com/google/uuid"
)

// User reprezentuje użytkownika aplikacji
type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // nigdy nie wysyłamy hasza w JSONie
	Location     string    `json:"location,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	AvatarURL    string    `json:"avatarUrl,omitempty"`
	Rating       float64   `json:"rating"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// UserRegister reprezentuje dane do rejestracji
type UserRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Location string `json:"location"`
}

// UserLogin reprezentuje dane do logowania
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserUpdate reprezentuje dane do aktualizacji profilu
type UserUpdate struct {
	Name      string `json:"name"`
	Location  string `json:"location"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatarUrl"`
}
