package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email,omitempty" validate:"required,email"`
	Password string    `json:"password,omitempty" validate:"required,min=6,max=12"`
}
