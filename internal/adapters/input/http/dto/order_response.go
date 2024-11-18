package dto

import "github.com/google/uuid"

type OrderResponse struct {
	UUID       uuid.UUID      `json:"id"`
	Status     string         `json:"status"`
	ExternalId string         `json:"external_id"`
	Amount     string         `json:"total_amount"`
	CreatedAt  string         `json:"created_at"`
	User       UserMeResponse `json:"user"`
	OrderItems []OrderItemResponse
}
