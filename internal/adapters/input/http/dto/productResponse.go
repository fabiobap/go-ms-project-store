package dto

import (
	"github.com/google/uuid"
)

type ProductResponse struct {
	Id          int64            `json:"id"`
	UUID        uuid.UUID        `json:"uuid"`
	Description string           `json:"description"`
	Amount      string           `json:"amount"`
	CategoryId  int64            `json:"category_id"`
	Image       string           `json:"image"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
	Category    CategoryResponse `json:"category"`
}
