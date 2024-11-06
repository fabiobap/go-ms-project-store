package dto

import (
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type ProductValidator interface {
	Validate() *helpers.ValidationResponse
}

type NewProductRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=250"`
	CategoryId  int64  `json:"category_id" validate:"required,min=1,max=9223372036854775807"`
	Description string `json:"description" validate:"required,min=3,max=25000"`
	Amount      int32  `json:"amount" validate:"required,min=1,max=9223372036854775807"`
	Slug        string `json:"slug" validate:"omitempty,min=3,max=250"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=250"`
	CategoryId  int64  `json:"category_id" validate:"required,min=1,max=9223372036854775807"`
	Description string `json:"description" validate:"required,min=3,max=25000"`
	Amount      int32  `json:"amount" validate:"required,min=1,max=9223372036854775807"`
	Slug        string `json:"slug" validate:"omitempty,min=3,max=250"`
}

func (ncr *NewProductRequest) Validate() *helpers.ValidationResponse {
	return helpers.ValidateRequests(ncr)
}

func (cur *UpdateProductRequest) Validate() *helpers.ValidationResponse {
	return helpers.ValidateRequests(cur)
}

// ValidateProduct is now a generic function that can handle any ProductValidator
func ValidateProduct(product ProductValidator) *helpers.ValidationResponse {
	return product.Validate()
}
