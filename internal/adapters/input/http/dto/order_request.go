package dto

import (
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type OrderValidator interface {
	Validate() *helpers.ValidationResponse
}

type CardRequest struct {
	Number   string `json:"number" validate:"required,min=13,max=16"`
	ExpMonth string `json:"exp_month" validate:"required,len=2"`
	ExpYear  string `json:"exp_year" validate:"required,len=4"`
	CVC      string `json:"cvc" validate:"required,len=3"`
	Name     string `json:"name" validate:"required"`
}

type ProductRequest struct {
	ID       string `json:"id" validate:"required,uuid4"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

type NewOrderRequest struct {
	Card     CardRequest      `json:"card" validate:"required"`
	Products []ProductRequest `json:"products" validate:"required,min=1,dive"`
}

func (ncr *NewOrderRequest) Validate() *helpers.ValidationResponse {
	return helpers.ValidateRequests(ncr)
}

func ValidateOrder(order OrderValidator) *helpers.ValidationResponse {
	return order.Validate()
}
