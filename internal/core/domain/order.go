package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/google/uuid"
)

type Order struct {
	ID         uint64    `db:"id"`
	UUID       uuid.UUID `db:"uuid"`
	ExternalId string    `db:"external_id"`
	Status     string    `db:"status"`
	Amount     int32     `db:"amount"`
	UserId     uint64    `db:"user_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	User       User
	OrderItems OrderItems
}

type Orders []Order

func (o Order) ToOrderDTO() dto.OrderResponse {
	amount := float64(o.Amount) / 100

	orderItems := make([]dto.OrderItemResponse, len(o.OrderItems))
	for i, orderItem := range o.OrderItems {
		orderItems[i] = orderItem.ToOrderItemDTO()
	}

	return dto.OrderResponse{
		UUID:       o.UUID,
		Status:     o.Status,
		ExternalId: o.ExternalId,
		Amount:     helpers.NumberFormat(amount, 2, ".", ","),
		CreatedAt:  helpers.DatetimeToString(o.CreatedAt),
		User: dto.UserMeResponse{
			ID:        o.User.UUID,
			Name:      o.User.Name,
			Email:     o.User.Email,
			CreatedAt: helpers.DatetimeToString(o.User.CreatedAt),
		},
		OrderItems: orderItems,
	}
}

func (o Orders) ToDTO() []dto.OrderResponse {
	dtos := make([]dto.OrderResponse, len(o))
	for i, order := range o {
		dtos[i] = order.ToOrderDTO()
	}
	return dtos
}
