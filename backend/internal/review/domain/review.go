package domain

import (
	"time"

	"github.com/google/uuid"
)

// Review reprezentuje ocenę/opinię o użytkowniku
type Review struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transactionId"`
	ReviewerID    uuid.UUID `json:"reviewerId"`
	Reviewer      *User     `json:"reviewer,omitempty"`
	ReviewedID    uuid.UUID `json:"reviewedId"`
	Reviewed      *User     `json:"reviewed,omitempty"`
	Rating        int       `json:"rating"`
	Comment       string    `json:"comment,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

// User w kontekście recenzji (uproszczony)
type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ReviewCreate reprezentuje dane do utworzenia nowej oceny
type ReviewCreate struct {
	TransactionID uuid.UUID `json:"transactionId" binding:"required"`
	ReviewedID    uuid.UUID `json:"reviewedId" binding:"required"`
	Rating        int       `json:"rating" binding:"required,min=1,max=5"`
	Comment       string    `json:"comment"`
}
