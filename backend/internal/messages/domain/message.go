package domain

import (
	"time"

	"github.com/google/uuid"
)

// Message reprezentuje wiadomość między użytkownikami
type Message struct {
	ID            uuid.UUID  `json:"id"`
	SenderID      uuid.UUID  `json:"senderId"`
	Sender        *User      `json:"sender,omitempty"`
	ReceiverID    uuid.UUID  `json:"receiverId"`
	Receiver      *User      `json:"receiver,omitempty"`
	TransactionID *uuid.UUID `json:"transactionId,omitempty"`
	Content       string     `json:"content"`
	Read          bool       `json:"read"`
	CreatedAt     time.Time  `json:"createdAt"`
}

// User w kontekście wiadomości (uproszczony)
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatarUrl,omitempty"`
}

// MessageCreate reprezentuje dane do utworzenia nowej wiadomości
type MessageCreate struct {
	ReceiverID    uuid.UUID  `json:"receiverId" binding:"required"`
	TransactionID *uuid.UUID `json:"transactionId"`
	Content       string     `json:"content" binding:"required"`
}

// ConversationPreview reprezentuje podgląd konwersacji w liście
type ConversationPreview struct {
	UserID        uuid.UUID `json:"userId"`
	UserName      string    `json:"userName"`
	UserAvatarURL string    `json:"userAvatarUrl,omitempty"`
	LastMessage   string    `json:"lastMessage"`
	LastMessageAt time.Time `json:"lastMessageAt"`
	UnreadCount   int       `json:"unreadCount"`
}
