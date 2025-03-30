package domain

import (
	"time"

	"github.com/google/uuid"
)

// Transaction reprezentuje transakcję wypożyczenia/wymiany książki
type Transaction struct {
	ID              uuid.UUID  `json:"id"`
	BookID          uuid.UUID  `json:"bookId"`
	Book            *Book      `json:"book,omitempty"`
	LenderID        uuid.UUID  `json:"lenderId"` // właściciel książki
	Lender          *User      `json:"lender,omitempty"`
	BorrowerID      uuid.UUID  `json:"borrowerId"` // wypożyczający
	Borrower        *User      `json:"borrower,omitempty"`
	Status          string     `json:"status"`
	TransactionType string     `json:"transactionType"` // lending, exchange
	StartDate       *time.Time `json:"startDate,omitempty"`
	DueDate         *time.Time `json:"dueDate,omitempty"`
	ReturnDate      *time.Time `json:"returnDate,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// User w kontekście transakcji (uproszczony)
type User struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Rating float64   `json:"rating"`
}

// Book w kontekście transakcji (uproszczony)
type Book struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ImageURLs []string  `json:"imageUrls,omitempty"`
}

// TransactionCreate reprezentuje dane do utworzenia nowej transakcji
type TransactionCreate struct {
	BookID          uuid.UUID  `json:"bookId" binding:"required"`
	BorrowerID      uuid.UUID  `json:"borrowerId" binding:"required"`
	TransactionType string     `json:"transactionType" binding:"required,oneof=lending exchange"`
	DueDate         *time.Time `json:"dueDate"`
	Notes           string     `json:"notes"`
}

// TransactionUpdate reprezentuje dane do aktualizacji transakcji
type TransactionUpdate struct {
	Status     string     `json:"status" binding:"omitempty,oneof=pending active completed canceled"`
	DueDate    *time.Time `json:"dueDate"`
	ReturnDate *time.Time `json:"returnDate"`
	Notes      string     `json:"notes"`
}
