package model

import "github.com/google/uuid"

type Chat struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Messages    Messages  `json:"messages"`
}
