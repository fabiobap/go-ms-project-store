package dto

import (
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type CategoryValidator interface {
	Validate() *helpers.ValidationResponse
}

type NewCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=250"`
	Slug string `json:"slug" validate:"omitempty,min=3,max=250"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=250"`
	Slug string `json:"slug" validate:"omitempty,min=3,max=250"`
}

func (ncr *NewCategoryRequest) Validate() *helpers.ValidationResponse {
	return helpers.ValidateRequests(ncr)
}

func (cur *UpdateCategoryRequest) Validate() *helpers.ValidationResponse {
	return helpers.ValidateRequests(cur)
}

// ValidateCategory is now a generic function that can handle any CategoryValidator
func ValidateCategory(category CategoryValidator) *helpers.ValidationResponse {
	return category.Validate()
}
