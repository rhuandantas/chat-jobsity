package model

import (
	"github.com/google/uuid"
)

type Message struct {
	ID     uuid.UUID `json:"id"`
	ChatID string    `param:"chatID,omitempty"`
	Text   string    `json:"text" validate:"required"`
	Author string    `json:"author" validate:"required"`
	Time   string    `json:"time,omitempty"`
}

type Messages []Message
