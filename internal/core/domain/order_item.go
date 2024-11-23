package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type OrderItem struct {
	ID        uint64    `db:"id"`
	Amount    int32     `db:"amount"`
	Quantity  int32     `db:"quantity"`
	OrderId   uint64    `db:"order_id"`
	ProductId uint64    `db:"product_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Product   Product
}

type OrderItems []OrderItem

func NewOrderItem(req dto.NewOrderItemDTO) OrderItem {
	return OrderItem{
		Amount:    req.Amount,
		Quantity:  req.Quantity,
		ProductId: req.ProductId,
		OrderId:   req.OrderId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (oi OrderItem) ToOrderItemDTO() dto.OrderItemResponse {
	amountOrderItem := float64(oi.Amount) / 100
	// amountProduct := float64(oi.Product.Amount) / 100

	return dto.OrderItemResponse{
		Amount:   helpers.NumberFormat(amountOrderItem, 2, ".", ","),
		Quantity: oi.Quantity,
		Product:  oi.Product.ToPublicProductDTO(),
	}
}

func (c OrderItems) ToDTO() []dto.OrderItemResponse {
	dtos := make([]dto.OrderItemResponse, len(c))
	for i, orderItem := range c {
		dtos[i] = orderItem.ToOrderItemDTO()
	}
	return dtos
}
