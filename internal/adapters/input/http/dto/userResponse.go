package dto

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	Id              int64     `json:"id"`
	UUID            uuid.UUID `json:"uuid"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	EmailVerifiedAt string    `json:"email_verified_at"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}
