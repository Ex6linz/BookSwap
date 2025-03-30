package domain

import (
	"time"

	"github.com/google/uuid"
)

// Book reprezentuje książkę w systemie
type Book struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description,omitempty"`
	ISBN        string    `json:"isbn,omitempty"`
	CategoryID  uuid.UUID `json:"categoryId"`
	Category    *Category `json:"category,omitempty"`
	Condition   string    `json:"condition"`
	OwnerID     uuid.UUID `json:"ownerId"`
	Owner       *User     `json:"owner,omitempty"`
	Status      string    `json:"status"`
	ImageURLs   []string  `json:"imageUrls,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// User w kontekście książki (uproszczony)
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Rating   float64   `json:"rating"`
	Location string    `json:"location,omitempty"`
}

// Category reprezentuje kategorię książki
type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

// BookImage reprezentuje zdjęcie książki
type BookImage struct {
	ID        uuid.UUID `json:"id"`
	BookID    uuid.UUID `json:"bookId"`
	ImageURL  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

// BookCreate reprezentuje dane do utworzenia nowej książki
type BookCreate struct {
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Description string    `json:"description"`
	ISBN        string    `json:"isbn"`
	CategoryID  uuid.UUID `json:"categoryId" binding:"required"`
	Condition   string    `json:"condition" binding:"required"`
	ImageBase64 []string  `json:"imageBase64,omitempty"`
}

// BookUpdate reprezentuje dane do aktualizacji książki
type BookUpdate struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	ISBN        string    `json:"isbn"`
	CategoryID  uuid.UUID `json:"categoryId"`
	Condition   string    `json:"condition"`
	Status      string    `json:"status"`
}

// BookFilter reprezentuje parametry filtrowania książek
type BookFilter struct {
	Title      string    `form:"title"`
	Author     string    `form:"author"`
	CategoryID uuid.UUID `form:"categoryId"`
	OwnerID    uuid.UUID `form:"ownerId"`
	Status     string    `form:"status"`
}
